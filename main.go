package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/jenkins-x/jx/pkg/kube"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"log"
	"net/http"
	"os"
	"time"

	jenkinsv1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"
	jenkinsclientv1 "github.com/jenkins-x/jx/pkg/client/clientset/versioned/typed/jenkins.io/v1"
	"k8s.io/client-go/rest"

	"github.com/jenkins-x/ext-spotbugs/findbugs"

	"github.com/pkg/errors"
)

// TODO replace these with imports from jx
const FactTypeStaticProgramAnalysis = "jx.staticProgramAnalysis"
const MeasurementCount   = "count"

func watch() (err error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	ns := os.Getenv("TEAM_NAMESPACE")
	if ns == "" {
		ns = "jx"
	}

	client, err := jenkinsclientv1.NewForConfig(config)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var httpClient = &http.Client{
		Transport: tr,
		Timeout: time.Second * 10,
	}

	listWatch := cache.NewListWatchFromClient(client.RESTClient(), "pipelineactivities", ns, fields.Everything())
	kube.SortListWatchByName(listWatch)
	_, actController := cache.NewInformer(
		listWatch,
		&jenkinsv1.PipelineActivity{},
		time.Minute*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onPipelineActivityObj(obj, httpClient, client)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				onPipelineActivityObj(newObj, httpClient, client)
			},
			DeleteFunc: func(obj interface{}) {

			},
		},
	)
	stop := make(chan struct{})
	go actController.Run(stop)

	return nil
}

func onPipelineActivityObj(obj interface{}, httpClient *http.Client, jxClient *jenkinsclientv1.JenkinsV1Client) {
	act, ok := obj.(*jenkinsv1.PipelineActivity)
	if !ok {
		log.Printf("unexpected type %s\n", obj)
	} else {
		err := onPipelineActivity(act, httpClient, jxClient)
		if err != nil {
			log.Println(err)
		}

	}
}

func onPipelineActivity(act *jenkinsv1.PipelineActivity, httpClient *http.Client, jxClient *jenkinsclientv1.JenkinsV1Client) (err error) {
	for _, attachment := range act.Spec.Attachments {
		if attachment.Name == "spotbugs" {
			// TODO Handle having multiple attachments properly
			for _, url := range attachment.URLs {
				url = fmt.Sprintf("%s?version=%d", url, time.Now().UnixNano()/int64(time.Millisecond))
				bugCollection, err := parseSpotBugsReport(url, httpClient)
				if err != nil {
					log.Println(errors.Wrap(err, fmt.Sprintf("Unable to retrieve %s for processing", url)))
				}
				fact := jenkinsv1.Fact{}
				if fact.Name == "" {
					fact.FactType = FactTypeStaticProgramAnalysis
					fact.Original = jenkinsv1.Original{
						URL:      url,
						MimeType: "application/xml",
						Tags: []string{
							"spotbugsXml.xml",
						},
					}
					fact.Tags = []string{
						"spotbugs",
					}
				}
				categories := make(map[string]map[string]int, 0)
				measurements := make([]jenkinsv1.Measurement, 0)
				for _, b := range bugCollection.BugInstance {
					category, ok := categories[b.Category]
					if !ok {
						category = make(map[string]int, 0)
					}
					switch b.Priority {
					case 1:
						category[jenkinsv1.StaticProgramAnalysisHighPriority]++
					case 2:
						category[jenkinsv1.StaticProgramAnalysisNormalPriority]++
					case 3:
						category[jenkinsv1.StaticProgramAnalysisLowPriority]++
					case 5:
						category[jenkinsv1.StaticProgramAnalysisIgnored]++
					}
					categories[b.Category] = category
				}
				for k, v := range categories {
					for l, w := range v {
						measurements = append(measurements, createMeasurement(k, l, w))
					}
				}
				measurements = append(measurements, createMeasurement("summary", jenkinsv1.StaticProgramAnalysisTotalBugs, bugCollection.FindBugsSummary.TotalBugs))
				measurements = append(measurements, createMeasurement("summary", jenkinsv1.StaticProgramAnalysisHighPriority, bugCollection.FindBugsSummary.HighPriority))
				measurements = append(measurements, createMeasurement("summary", jenkinsv1.StaticProgramAnalysisNormalPriority, bugCollection.FindBugsSummary.NormalPriority))
				measurements = append(measurements, createMeasurement("summary", jenkinsv1.StaticProgramAnalysisLowPriority, bugCollection.FindBugsSummary.LowPriority))
				measurements = append(measurements, createMeasurement("summary", jenkinsv1.StaticProgramAnalysisIgnored, bugCollection.FindBugsSummary.IgnorePriority))
				measurements = append(measurements, createMeasurement("summary", jenkinsv1.StaticProgramAnalysisTotalClasses, bugCollection.FindBugsSummary.TotalClasses))
				fact.Measurements = measurements
				found := 0
				for i, f := range act.Spec.Facts {
					if f.FactType == jenkinsv1.FactTypeStaticProgramAnalysis {
						act.Spec.Facts[i] = fact
						found++
					}
				}
				if found > 1 {
					return errors.New(fmt.Sprintf("More than one fact of kind %s, found %d", FactTypeStaticProgramAnalysis, found))
				} else if found == 0 {
					act.Spec.Facts = append(act.Spec.Facts, fact)
				}
				act, err = jxClient.PipelineActivities(act.Namespace).Update(act)
				log.Printf("Updated PipelineActivity %s with data from %s\n", act.Name, url)
				if err != nil {
					log.Println(errors.Wrap(err, fmt.Sprintf("Error updating PipelineActivity %s", act.Name)))
				}
			}
		}
	}
	return nil
}

func parseSpotBugsReport(url string, httpClient *http.Client) (collection findbugs.BugCollection, err error) {
	response, err := httpClient.Get(url)
	if err != nil {
		return findbugs.BugCollection{}, err
	}
	if response.StatusCode > 299 || response.StatusCode < 200 {
		return findbugs.BugCollection{}, errors.New(fmt.Sprintf("Status code: %d, error: %s", response.StatusCode, response.Status))
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return findbugs.BugCollection{}, err
	}
	err = xml.Unmarshal(body, &collection)
	if err != nil {
		return findbugs.BugCollection{}, err
	}
	return collection, nil
}

func createMeasurement(t string, measurement string, value int) jenkinsv1.Measurement {
	return jenkinsv1.Measurement{
		Name:             fmt.Sprintf("%s-%s", t, measurement),
		MeasurementType:  MeasurementCount,
		MeasurementValue: value,
	}
}

func main() {
	err := watch()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := "Ready :-D"

	from := ""
	if r.URL != nil {
		from = r.URL.String()
	}
	if from != "/favicon.ico" {
		log.Printf("title: %s\n", title)
	}

	fmt.Fprintf(w, title+"\n")
}