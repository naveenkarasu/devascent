using System.Collections.Generic;

public class Solution {
    public object[] time_map_ops(object[][] operations) {
        var times = new Dictionary<string, List<long>>();
        var vals = new Dictionary<string, List<string>>();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            var op = operations[i];
            string name = (string)op[0];
            if (name == "set") {
                string key = (string)op[1];
                string value = (string)op[2];
                long ts = (long)op[3];
                if (!times.ContainsKey(key)) { times[key] = new List<long>(); vals[key] = new List<string>(); }
                times[key].Add(ts);
                vals[key].Add(value);
                out_[i] = null;
            } else { // get
                string key = (string)op[1];
                long ts = (long)op[2];
                if (!times.ContainsKey(key)) {
                    out_[i] = "";
                } else {
                    var tsList = times[key];
                    // bisect_right: find first index where tsList[mid] > ts
                    int lo = 0, hi = tsList.Count;
                    while (lo < hi) {
                        int mid = (lo + hi) / 2;
                        if (tsList[mid] <= ts) lo = mid + 1;
                        else hi = mid;
                    }
                    out_[i] = lo > 0 ? vals[key][lo - 1] : "";
                }
            }
        }
        return out_;
    }
}
