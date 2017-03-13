package lucene

import (
	"fmt"
	"os"
	"github.com/balzaczyy/golucene/core/store"
	"github.com/balzaczyy/golucene/core/document"
	"github.com/balzaczyy/golucene/core/util"
	"github.com/balzaczyy/golucene/core/index"
	"github.com/balzaczyy/golucene/core/search"
	std "github.com/balzaczyy/golucene/analysis/standard"
	"strings"
)

const DIR_INDEX string = "data_lucene"

type Call struct {
	ID   string `json:"id"`
	Url  string `json:"url"`
	Desc string `json:"desc"`
	Tags string `json:"tags"`
}

func Write(call Call) {

	directory, _ := store.OpenFSDirectory(DIR_INDEX)
	analyzer := std.NewStandardAnalyzer()
	conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer)
	writer, _ := index.NewIndexWriter(directory, conf)
	d := document.NewDocument()
	stripUrl := strings.Replace(call.Url, "\\", " ", -1)
	stripUrl = strings.Replace(stripUrl, "/", " ", -1)
	stripUrl = strings.Replace(stripUrl, ".", " ", -1)
	stripUrl = strings.Replace(stripUrl, ":", " ", -1)
	stripUrl = strings.Replace(stripUrl, ",", " ", -1)
	stripUrl = strings.Replace(stripUrl, "=", " ", -1)
	stripUrl = strings.Replace(stripUrl, "-", " ", -1)
	stripUrl = strings.Replace(stripUrl, "?", " ", -1)
	fmt.Println(call.Tags)
	fmt.Println(stripUrl)

	d.Add(document.NewTextFieldFromString("id", call.ID, document.STORE_YES))
	d.Add(document.NewTextFieldFromString("desc", call.Desc, document.STORE_YES))
	d.Add(document.NewTextFieldFromString("tags", call.Tags, document.STORE_YES))
	d.Add(document.NewTextFieldFromString("url", stripUrl, document.STORE_YES))
	d.Add(document.NewTextFieldFromString("all", call.ID+" "+stripUrl+" "+call.Tags+" "+call.Desc, document.STORE_YES))
	writer.AddDocument(d.Fields())
	writer.Close()
}

func Find(field string, query string) (map[string]string) {
	directory, _ := store.OpenFSDirectory(DIR_INDEX)
	reader, _ := index.OpenDirectoryReader(directory)
	searcher := search.NewIndexSearcher(reader)

	q := search.NewTermQuery(index.NewTerm(field, query))
	res, _ := searcher.Search(q, nil, 1000)
	results := make(map[string]string)
	for _, hit := range res.ScoreDocs {
		doc, _ := reader.Document(hit.Doc);
		id := doc.Get("id")
		all := doc.Get("all")
		//fmt.Println("found id " + id)
		//fmt.Println("found all " + all)
		results[id] = all;
	}
	return results
}

func Test() {
	util.SetDefaultInfoStream(util.NewPrintStreamInfoStream(os.Stdout))
	index.DefaultSimilarity = func() index.Similarity {
		return search.NewDefaultSimilarity()
	}

	directory, _ := store.OpenFSDirectory("test_index")
	analyzer := std.NewStandardAnalyzer()
	conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer)
	writer, _ := index.NewIndexWriter(directory, conf)

	d := document.NewDocument()
	d.Add(document.NewTextFieldFromString("foo", "bar", document.STORE_YES))
	writer.AddDocument(d.Fields())
	writer.Close() // ensure index is written

	reader, _ := index.OpenDirectoryReader(directory)
	searcher := search.NewIndexSearcher(reader)

	q := search.NewTermQuery(index.NewTerm("foo", "bar"))
	res, _ := searcher.Search(q, nil, 1000)
	fmt.Printf("Found %v hit(s).\n", res.TotalHits)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("foo -> %v\n", doc.Get("foo"))
	}

}
