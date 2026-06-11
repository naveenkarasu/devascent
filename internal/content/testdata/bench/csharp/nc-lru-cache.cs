using System.Collections.Generic;

public class Solution {
    public object[] lru_cache_ops(long capacity, object[][] operations) {
        var d = new LinkedList<(long key, long val)>();
        var map = new Dictionary<long, LinkedListNode<(long key, long val)>>();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            var op = operations[i];
            string name = (string)op[0];
            if (name == "get") {
                long key = (long)op[1];
                if (!map.ContainsKey(key)) {
                    out_[i] = -1L;
                } else {
                    var node = map[key];
                    long val = node.Value.val;
                    d.Remove(node);
                    map[key] = d.AddLast((key, val));
                    out_[i] = val;
                }
            } else { // put
                long key = (long)op[1];
                long val = (long)op[2];
                if (map.ContainsKey(key)) {
                    d.Remove(map[key]);
                }
                map[key] = d.AddLast((key, val));
                if (d.Count > capacity) {
                    var first = d.First!;
                    map.Remove(first.Value.key);
                    d.RemoveFirst();
                }
                out_[i] = null;
            }
        }
        return out_;
    }
}
