package main

import (
	"bleausky/db"
	"bleausky/rules"
	"bleausky/sync"
	"bytes"
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/bluesky-social/indigo/repo"
	"github.com/gorilla/websocket"
	"regexp"

	//"regexp"
	"slices"

	//"github.com/ipfs/go-cid"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// Firehose - mutation list? Sends a stream of (repo=user, cid=commit id, seq=sequence number)
// Then you Follow just some users to get a filtered list

// repo resolver
// repo syncer? (do I need to resolve to sync?)

func main() {
	//db.DropTables()
	if true {
		HydrateDBFromFirehose()
	} else if false {
		db.ListPosts()
		//db.CountPosts()
	} else if false {
		repo := "did%3Aplc%3A3gn2axv5p2vuyaijjog7apde"
		//repo := "did:plc:3gn2axv5p2vuyaijjog7apde"
		//repo := "3gn2axv5p2vuyaijjog7apde"
		sync.SyncRepo(repo, "0")
	} else if true {
		sync.GetPosts()
	}
}

func HydrateDBFromFirehose() {
	db.CreateDB()
	uri := "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
	con, _, err := websocket.DefaultDialer.Dial(uri, http.Header{})
	if err != nil {
		panic(err)
	}
	rsc := &events.RepoStreamCallbacks{
		RepoCommit: func(evt *atproto.SyncSubscribeRepos_Commit) error {
			rr, err := repo.ReadRepoFromCar(context.TODO(), bytes.NewReader(evt.Blocks))
			if err != nil {
				return err
			}

			for _, op := range evt.Ops {
				// if op.Path ends with "post"
				// then insert into db
				if strings.HasPrefix(op.Path, "app.bsky.feed.postgate") {
					// do nothing, wtf is a postgate?
				} else if strings.HasPrefix(op.Path, "app.bsky.feed.post") {
					// evt.Blobs is a CAR file diff of repo state
					if op.Action == "create" {
						//fmt.Println("Event from ", evt.Repo)
						//fmt.Printf(" - %s record %s\n", op.Action, op.Path)
						//db.InsertPost(evt.Repo, op.Cid.String(), evt.Seq)
						db.InsertPost(evt.Repo, op.Path, evt.Seq)

						_, rec, err := rr.GetRecord(context.TODO(), op.Path)
						if err != nil {
							return err
						}

						//fmt.Println("Record: ", rc, rec)
						//fmt.Println(rec)
						switch recV := rec.(type) {
						default:
							fmt.Println("default: ", recV)
						case *bsky.FeedPost:
							if recV.Reply == nil {
								if slices.Contains(recV.Langs, "en") {
									txt := recV.Text
									// replace all newlines with space in txt
									//txt = strings.ReplaceAll(txt, "\n", " ")
									spaceRe := regexp.MustCompile(`\s+`)
									txt = strings.TrimSpace(spaceRe.ReplaceAllString(txt, " "))
									if rules.PostFilter(txt) {
										fmt.Println(txt)
									} else {
										//fmt.Println("Filtered: ", txt)
									}
									//// if txt contains any non-whitespace characters
									//if strings.TrimSpace(txt) != "" {
									//	// text starts with a capital letter, and ends with a period, question mark, or exclamation mark
									//	if txt[0] >= 'A' && txt[0] <= 'Z' && (txt[len(txt)-1] == '.' || txt[len(txt)-1] == '?' || txt[len(txt)-1] == '!') {
									//		// contains any emoji
									//		if strings.ContainsAny(txt, "😀😃😄😁😆😅😂🤣😭😢😥😰😓😩😫😨😱😠😡😤😖😆😋😷😎😴😵😲😟😦😧😈👿😮😬😐😕😯😶😇😏😑😒🙄🤔😳😞😟😠😡😔😕🙁☹️😣😖😫😩😤😠😡😶😐😑😯😦😧😮😲😵😳😱😨😰😢😥😓😭😂😅😆😩😫😓😢😥😰😭") {
									//			//fmt.Println("! ", txt)
									//		} else {
									//			fmt.Println("> ", txt)
									//		}
									//	}
									//}
								}
							}
						}
					}
				} else if strings.HasPrefix(op.Path, "app.bsky.feed.repost") {
					// do nothing
				} else if strings.HasPrefix(op.Path, "app.bsky.feed.like") {
					// do nothing
				} else if strings.HasPrefix(op.Path, "app.bsky.graph.follow") {
					// do nothing
				} else {
				}
			}
			return nil
		},
	}
	sched := sequential.NewScheduler("myfirehose", rsc.EventHandler)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	events.HandleRepoStream(context.Background(), con, sched, logger)
}
