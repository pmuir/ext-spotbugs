package findbugs

import "encoding/xml"

// The FindBugs XML model
type BugCollection struct {
	XMLName           xml.Name        `xml:"BugCollection"`
	Sequence          int             `xml:"sequence,attr"`
	Release           string          `xml:"release,attr"`
	AnalysisTimestamp string          `xml:"analysisTimestamp,attr"`
	Version           string          `xml:"version,attr"`
	Timestamp         string          `xml:"timestamp,attr"`
	Projects          Project         `xml:"Project"`
	Errors            Errors          `xml:"Errors"`
	FindBugsSummary   FindBugsSummary `xml:"FindBugsSummary"`
	BugInstance       []BugInstance   `xml:"BugInstance"`
	BugCategory       []BugCategory   `xml:"BugCategory"`
	BugPattern        []BugPattern    `xml:"BugPattern"`
	BugCode           []BugCode       `xml:"BugCode"`
}

type Project struct {
	XMLName           xml.Name `xml:"Project"`
	ProjectName       string   `xml:"projectName,attr"`
	Jar               []string `xml:"Jar"`
	AuxClasspathEntry []string `xml:"AuxClasspathEntry"`
	SrcDir            []string `xml:"SrcDir"`
	WrkDir            string   `xml:"WrkDir"`
	Errors            Errors   `xml:"Errors"`
}

type Errors struct {
	XMLName        xml.Name `xml:"Errors"`
	MissingClasses int      `xml:"missingClasses,attr"`
	Errors         int      `xml:"errors,attr"`
}

type FindBugsSummary struct {
	XMLName           xml.Name        `xml:"FindBugsSummary"`
	NumPackages       int             `xml:"num_packages,attr"`
	TotalClasses      int             `xml:"total_classes,attr"`
	HighPriority      int             `xml:"priority_1,attr"`
	NormalPriority    int             `xml:"priority_2,attr"`
	LowPriority       int             `xml:"priority_3,attr"`
	IgnorePriority    int             `xml:"priority_5,attr"`
	ExpPriority       int             `xml:"priority_4,attr"`
	TotalSize         int             `xml:"total_size,attr"`
	ClockSeconds      float32         `xml:"clock_seconds,attr"`
	ReferencedClasses int             `xml:"referenced_classes,attr"`
	VMVersion         string          `xml:"vm_version,attr"`
	TotalBugs         int             `xml:"total_bugs,attr"`
	JavaVersion       string          `xml:"java_version,attr"`
	GCSeconds         float32         `xml:"gc_seconds,attr"`
	AllocMBytes       float32         `xml:"alloc_mbytes,attr"`
	CPUSeconds        float32         `xml:"cpu_seconds,attr"`
	PeakMBytes        float32         `xml:"peak_mbytes,attr"`
	Timestamp         string          `xml:"timestamp,attr"`
	PackageStats      []PackageStats  `xml:"PackageStats"`
	FindBugsProfile   FindBugsProfile `xml:"FindBugsProfile"`
	ClassFeatures     ClassFeatures   `xml:"ClassFeatures"`
	History           History         `xml:"History"`
}

type PackageStats struct {
	XMLName    xml.Name     `xml:"PackageStats"`
	TotalBugs  int          `xml:"total_bugs,attr"`
	TotalSize  int          `xml:"total_size,attr"`
	TotalTypes int          `xml:"total_types,attr"`
	Package    string       `xml:"package,attr"`
	ClassStats []ClassStats `xml:"ClassStats"`
}

type ClassStats struct {
	XMLName    xml.Name `xml:"ClassStats"`
	Bugs       int      `xml:"bugs,attr"`
	Size       int      `xml:"size,attr"`
	Interface  bool     `xml:"interface,attr"`
	SourceFile string   `xml:"sourceFile,attr"`
	Class      string   `xml:"class,attr"`
}

type FindBugsProfile struct {
	XMLName      xml.Name       `xml:"FindBugsProfile"`
	ClassProfile []ClassProfile `xml:"ClassProfile"`
}

type ClassProfile struct {
	XMLName                                    xml.Name `xml:"ClassProfile"`
	AvgMicrosecondsPerInvocation               int      `xml:"avgMicrosecondsPerInvocation,attr"`
	TotalMicrosecondsPerInvocation             int      `xml:"totalMicrosecondsPerInvocation,attr"`
	StandardDeviationMicrosecondsPerInvocation int      `xml:"standardDeviationMicrosecondsPerInvocation,attr"`
	TotalMilliseconds                          int      `xml:"totalMilliseconds,attr"`
	Invocations                                int      `xml:"invocations,attr"`
	Name                                       string   `xml:"nameFile,attr"`
}

type ClassFeatures struct {
	XMLName xml.Name `xml:"ClassFeatures"`
}

type History struct {
	XMLName xml.Name `xml:"History"`
}

type BugInstance struct {
	XMLName               xml.Name   `xml:"BugInstance"`
	InstanceOccurenceNum  int        `xml:"instanceOccurenceNum,attr"`
	InstanceHash          string     `xml:"instanceHash,attr"`
	Cweid                 int        `xml:"cweid,attr"`
	Rank                  int        `xml:"rank,attr"`
	Abbrev                string     `xml:"abbrev,attr"`
	Category              string     `xml:"category,attr"`
	Priority              int        `xml:"priority,attr"`
	Type                  string     `xml:"type,attr"`
	InstanceOccurrenceMax int        `xml:"instanceOccurenceMax,attr"`
	ShortMessage          string     `xml:"ShortMessage"`
	LongMessage           string     `xml:"LongMessage"`
	Class                 Class      `xml:"Class"`
	Method                Method     `xml:"Method"`
	SourceLine            SourceLine `xml:"SourceLine"`
	Field                 Field      `xml:"Field"`
}

type Class struct {
	XMLName    xml.Name   `xml:"Class"`
	ClassName  string     `xml:"classname, attr"`
	Primary    bool       `xml:"primary, attr"`
	SourceLine SourceLine `xml:"SourceLine"`
	Message    string     `xml:"Message"`
}

type SourceLine struct {
	XMLName       xml.Name `xml:"SourceLine"`
	ClassName     string   `xml:"classname, attr"`
	SourcePath    string   `xml:"sourcepath, attr"`
	SourceFile    string   `xml:"sourcefile, attr"`
	Start         int      `xml:"start, attr"`
	End           int      `xml:"end, attr"`
	StartBytecode int      `xml:"startBytecode, attr"`
	EndBytecode   int      `xml:"endBytecode, attr"`
	Message       string   `xml:"Message"`
}

type Method struct {
	XMLName    xml.Name   `xml:"Method"`
	IsStatic   bool       `xml:"isStatic, attr"`
	ClassName  string     `xml:"classname, attr"`
	Signature  string     `xml:"signature, attr"`
	Name       string     `xml:"name, attr"`
	Primary    bool       `xml:"primary, attr"`
	SourceLine SourceLine `xml:"SourceLine"`
	Message    string     `xml:"Message"`
}

type Field struct {
	XMLName    xml.Name   `xml:"Field"`
	IsStatic   bool       `xml:"isStatic, attr"`
	ClassName  string     `xml:"classname, attr"`
	Signature  string     `xml:"signature, attr"`
	Name       string     `xml:"name, attr"`
	Primary    bool       `xml:"primary, attr"`
	SourceLine SourceLine `xml:"SourceLine"`
	Message    string     `xml:"Message"`
}

type BugCategory struct {
	XMLName  xml.Name `xml:"BugCategory"`
	Category string   `xml:"category, attr"`
	Message  string   `xml:"Message"`
}

type BugPattern struct {
	XMLName          xml.Name `xml:"BugPattern"`
	Abbrev           string   `xml:"abbrev, attr"`
	Type             string   `xml:"type, attr"`
	Category         string   `xml:"category, attr"`
	ShortDescription string   `xml:"ShortDescription"`
	Details          string   `xml:"Details"`
}

type BugCode struct {
	XMLName     xml.Name `xml:"BugCode"`
	Cweid       int      `xml:"cweid,attr"`
	Abbrev      string   `xml:"abbrev,attr"`
	Description string   `xml:"Description"`
}
