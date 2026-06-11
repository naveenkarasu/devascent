using System.Collections.Generic;

public class Solution {
    public object[] detect_squares_ops(object[][] operations) {
        var cnt = new Dictionary<string, long>();
        var points = new List<long[]>();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            var op = operations[i];
            string name = (string)op[0];
            var pt = (object[])op[1];
            long px = (long)pt[0];
            long py = (long)pt[1];
            if (name == "add") {
                string key = px + "," + py;
                if (!cnt.ContainsKey(key)) points.Add(new long[] { px, py });
                cnt[key] = cnt.GetValueOrDefault(key, 0L) + 1L;
                out_[i] = null;
            } else { // count
                long total = 0;
                foreach (var p in points) {
                    long x = p[0], y = p[1];
                    if (Math.Abs(x - px) == Math.Abs(y - py) && x != px && y != py) {
                        long c = cnt[x + "," + y];
                        long c1 = cnt.GetValueOrDefault(px + "," + y, 0L);
                        long c2 = cnt.GetValueOrDefault(x + "," + py, 0L);
                        total += c * c1 * c2;
                    }
                }
                out_[i] = total;
            }
        }
        return out_;
    }
}
