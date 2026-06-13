using System.Collections.Generic;
using System.Linq;

public class Solution {
    public long[] set_union(long[] a, long[] b) {
        var seen = new HashSet<long>(a);
        foreach (long x in b) seen.Add(x);
        var result = seen.ToList();
        result.Sort();
        return result.ToArray();
    }
}
