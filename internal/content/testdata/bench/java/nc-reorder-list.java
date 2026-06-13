class Solution {
    public ListNode reorder_list(ListNode head) {
        if (head == null || head.next == null) return head;
        // Find middle
        ListNode slow = head, fast = head;
        while (fast.next != null && fast.next.next != null) {
            slow = slow.next;
            fast = fast.next.next;
        }
        // Reverse second half
        ListNode second = slow.next;
        slow.next = null;
        ListNode prev = null;
        while (second != null) {
            ListNode nxt = second.next;
            second.next = prev;
            prev = second;
            second = nxt;
        }
        // Merge two halves
        ListNode first = head;
        while (prev != null) {
            ListNode n1 = first.next;
            ListNode n2 = prev.next;
            first.next = prev;
            prev.next = n1;
            first = n1;
            prev = n2;
        }
        return head;
    }
}
