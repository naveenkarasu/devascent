use std::collections::HashMap;
use std::collections::HashSet;

fn twitter_ops(operations: Vec<Vec<J>>) -> Vec<J> {
    let mut time: i64 = 0;
    let mut tweets: HashMap<i64, Vec<(i64, i64)>> = HashMap::new(); // user -> (time, tweetId)
    let mut following: HashMap<i64, HashSet<i64>> = HashMap::new();
    let mut out: Vec<J> = Vec::new();

    for op in &operations {
        let name = match &op[0] {
            J::Str(s) => s.as_str(),
            _ => "",
        };
        match name {
            "postTweet" => {
                let u = if let J::Int(v) = &op[1] { *v } else { 0 };
                let t = if let J::Int(v) = &op[2] { *v } else { 0 };
                tweets.entry(u).or_insert_with(Vec::new).push((time, t));
                time += 1;
                out.push(J::Null);
            }
            "getNewsFeed" => {
                let u = if let J::Int(v) = &op[1] { *v } else { 0 };
                let mut people: HashSet<i64> = HashSet::new();
                if let Some(f) = following.get(&u) {
                    for &p in f {
                        people.insert(p);
                    }
                }
                people.insert(u);
                let mut all_tweets: Vec<(i64, i64)> = Vec::new();
                for &p in &people {
                    if let Some(ts) = tweets.get(&p) {
                        for &t in ts {
                            all_tweets.push(t);
                        }
                    }
                }
                // sort by (time, tweetId) descending
                all_tweets.sort();
                all_tweets.reverse();
                let mut feed: Vec<J> = Vec::new();
                for &(_, tid) in all_tweets.iter().take(10) {
                    feed.push(J::Int(tid));
                }
                out.push(J::Arr(feed));
            }
            "follow" => {
                let a = if let J::Int(v) = &op[1] { *v } else { 0 };
                let b = if let J::Int(v) = &op[2] { *v } else { 0 };
                following.entry(a).or_insert_with(HashSet::new).insert(b);
                out.push(J::Null);
            }
            _ => {
                // unfollow
                let a = if let J::Int(v) = &op[1] { *v } else { 0 };
                let b = if let J::Int(v) = &op[2] { *v } else { 0 };
                if let Some(f) = following.get_mut(&a) {
                    f.remove(&b);
                }
                out.push(J::Null);
            }
        }
    }
    out
}
