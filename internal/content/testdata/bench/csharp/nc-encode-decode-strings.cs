public class Solution {
    public string[] encode_decode(string[] strs) {
        // Encode
        var sb = new System.Text.StringBuilder();
        foreach (var s in strs) {
            sb.Append(s.Length);
            sb.Append('#');
            sb.Append(s);
        }
        string encoded = sb.ToString();

        // Decode
        var res = new List<string>();
        int i = 0;
        while (i < encoded.Length) {
            int j = i;
            while (encoded[j] != '#') j++;
            int length = int.Parse(encoded.Substring(i, j - i));
            res.Add(encoded.Substring(j + 1, length));
            i = j + 1 + length;
        }
        return res.ToArray();
    }
}
