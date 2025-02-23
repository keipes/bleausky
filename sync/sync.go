package sync

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/data"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/ipfs/go-cid"
	"io"
	"net/http"
	"os"
)

func SyncRepo(repoId string, revision string) {
	// 	https://bsky.network/xrpc/com.atproto.sync.getRepo?did=did%3Aplc%3A3gn2axv5p2vuyaijjog7apde
	//host := "morel.us-east.host.bsky.network"
	//host := "bsky.network"
	//host := "bsky.social"
	host := "helvella.us-east.host.bsky.network"
	endpoint := "com.atproto.sync.getRepo"
	uri := fmt.Sprintf("https://%s/xrpc/%s?did=%s&since=%s", host, endpoint, repoId, revision)
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//bodyStr := string(body)
	decodedCarFile, err := repo.ReadRepoFromCar(context.TODO(), resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(decodedCarFile)
	decodedCarFile.ForEach(context.TODO(), "", func(k string, v cid.Cid) error {
		fmt.Println(k, v)
		return nil
	})
	// resp.Body is an io.ReadCloser, convert the body to a string
	//body := resp.Body
	// response to string then print
	//fmt.Println(bodyStr)
}

func GetPosts() {
	post := "app.bsky.feed.post%2F3lgxzmtjdss2j"
	repoId := "did%3Aplc%3A3gn2axv5p2vuyaijjog7apde"
	atUri := fmt.Sprintf("at://%s/%s", repoId, post)
	//host := "helvella.us-east.host.bsky.network"
	host := "public.api.bsky.app."
	endpoint := "app.bsky.feed.getPosts"
	uri := fmt.Sprintf("https://%s/xrpc/%s?uris=%s", host, endpoint, atUri)
	if true {
		resp, err := http.Get(uri)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		dataMap, err := data.UnmarshalJSON(body)
		posts := dataMap["posts"]
		for _, post := range posts.([]interface{}) {
			// iterate and print key value pairs from each post
			for k, v := range post.(map[string]interface{}) {
				if k == "record" {
					record := v.(map[string]interface{})
					for k, v := range record {
						//fmt.Printf("%s: %v\n", k, v)
						if k == "reply" {
							reply := v.(map[string]interface{})
							for k, v := range reply {
								fmt.Printf("%s: %v\n", k, v)
							}
						}
					}
					//continue
				}
				//fmt.Printf("%s: %v\n", k, v)
			}
			//fmt.Println(post)
		}
	} else {
		//client := util.RobustHTTPClient()
		ctx := context.TODO()
		//xrpcc, err := cliutil.GetXrpcClient(ctx, true)
		client, err := MakeXRPCClient(ctx)
		if err != nil {
			fmt.Println(err)
		}
		bsky.FeedGetPosts(context.TODO(), client, []string{atUri})
	}

	//bodyStr := string(body)
	//fmt.Println(bodyStr)
	//bytes, err := json.MarshalIndent(bodyStr, "", "  ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(bytes))
}

func MakeXRPCClient(ctx context.Context) (*xrpc.Client, error) {
	username := os.Getenv("BSKY_USERNAME")
	pass := os.Getenv("BSKY_PASSWORD")

	xrpcc := &xrpc.Client{
		//Host: "https://bsky.social",
		Host: "https://public.api.bsky.app",
		Auth: &xrpc.AuthInfo{Handle: username},
	}

	auth, err := atproto.ServerCreateSession(ctx, xrpcc, &atproto.ServerCreateSession_Input{
		Identifier: xrpcc.Auth.Handle,
		Password:   pass,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	xrpcc.Auth.Did = auth.Did
	xrpcc.Auth.AccessJwt = auth.AccessJwt
	xrpcc.Auth.RefreshJwt = auth.RefreshJwt

	return xrpcc, nil
}
