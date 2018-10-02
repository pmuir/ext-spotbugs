package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jenkinsv1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"
	jenkinsclientv1 "github.com/jenkins-x/jx/pkg/client/clientset/versioned/typed/jenkins.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/jenkins-x/jenkins-x-spotbugs-reporter/findbugs"

	"github.com/pkg/errors"
)

func watch() (err error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	ns := os.Getenv("SPOTBUGS_NAMESPACE")
	client, err := jenkinsclientv1.NewForConfig(config)
	if err != nil {
		return err
	}
	watch, err := client.PipelineActivities(ns).Watch(metav1.ListOptions{})
	if err != nil {
		return err
	}

	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	for event := range watch.ResultChan() {
		act, ok := event.Object.(*jenkinsv1.PipelineActivity)
		if !ok {
			log.Fatalf("unexpected type %s\n", event)
		}
		//
		if act.Spec.Summaries.StaticProgramAnalysis.TotalClasses == 0 {
			for _, attachment := range act.Spec.Attachments {
				if attachment.Name == "spotbugs" {
					// TODO Handle having multiple attachments properly
					for _, url := range attachment.URLs {
						url = fmt.Sprintf("%s?version=%d", url, time.Now().UnixNano()/int64(time.Millisecond))
						bugCollection, err := parseSpotBugsReport(url, httpClient)
						if err != nil {
							log.Println(errors.Wrap(err, fmt.Sprintf("Unable to retrieve %s for processing", url)))
						}
						// Create the summaries for the categories
						categories := make(map[string]jenkinsv1.StaticProgramAnalysisCategory)
						for _, b := range bugCollection.BugInstance {
							category, ok := categories[b.Category]
							if !ok {
								category = jenkinsv1.StaticProgramAnalysisCategory{}
							}
							switch b.Priority {
							case 1:
								category.HighPriority++
							case 2:
								category.NormalPriority++
							case 3:
								category.LowPriority++
							case 5:
								category.Ignored++
							}
							categories[b.Category] = category
						}
						act.Spec.Summaries.StaticProgramAnalysis = jenkinsv1.StaticProgramAnalysis{
							TotalBugs:      bugCollection.FindBugsSummary.TotalBugs,
							HighPriority:   bugCollection.FindBugsSummary.HighPriority,
							NormalPriority: bugCollection.FindBugsSummary.NormalPriority,
							LowPriority:    bugCollection.FindBugsSummary.LowPriority,
							Ignored:        bugCollection.FindBugsSummary.IgnorePriority,
							TotalClasses:   bugCollection.FindBugsSummary.TotalClasses,
							//Categories:     categories,
						}
						act, err = client.PipelineActivities(act.Namespace).Update(act)
						log.Printf("Updated PipelineActivity %s with data from %s\n", act.Name, url)
						if err != nil {
							log.Println(errors.Wrap(err, fmt.Sprintf("Error updating PipelineActivity %s", act.Name)))
						}
					}
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

func main() {
	err := watch()
	if err != nil {
		panic(err.Error())
	}

}
