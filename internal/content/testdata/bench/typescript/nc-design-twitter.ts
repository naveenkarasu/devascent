class Twitter {
    private time: number;
    private tweets: Map<number, [number, number][]>;
    private following: Map<number, Set<number>>;
    constructor() {
        this.time = 0;
        this.tweets = new Map();
        this.following = new Map();
    }
    post_tweet(user: number, tweet: number): void {
        if (!this.tweets.has(user)) this.tweets.set(user, []);
        this.tweets.get(user)!.push([this.time, tweet]);
        this.time++;
    }
    get_news_feed(user: number): number[] {
        const people = new Set(this.following.get(user) || []);
        people.add(user);
        const allTweets: [number, number][] = [];
        for (const p of people) {
            const tw = this.tweets.get(p) || [];
            for (const t of tw) allTweets.push(t);
        }
        allTweets.sort((a, b) => b[0] - a[0] || b[1] - a[1]);
        return allTweets.slice(0, 10).map(t => t[1]);
    }
    follow(a: number, b: number): void {
        if (!this.following.has(a)) this.following.set(a, new Set());
        this.following.get(a)!.add(b);
    }
    unfollow(a: number, b: number): void {
        if (this.following.has(a)) this.following.get(a)!.delete(b);
    }
}

function twitter_ops(operations: any[][]): any[] {
    const tw = new Twitter();
    const out: any[] = [];
    for (const op of operations) {
        if (op[0] === "postTweet") {
            tw.post_tweet(op[1], op[2]);
            out.push(null);
        } else if (op[0] === "getNewsFeed") {
            out.push(tw.get_news_feed(op[1]));
        } else if (op[0] === "follow") {
            tw.follow(op[1], op[2]);
            out.push(null);
        } else {
            tw.unfollow(op[1], op[2]);
            out.push(null);
        }
    }
    return out;
}
