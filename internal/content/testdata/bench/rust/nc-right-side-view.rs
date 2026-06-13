use std::rc::Rc;
use std::cell::RefCell;
use std::collections::VecDeque;

fn right_side_view(root: Option<Rc<RefCell<TreeNode>>>) -> Vec<i64> {
    let mut res: Vec<i64> = Vec::new();
    if root.is_none() {
        return res;
    }
    let mut q: VecDeque<Rc<RefCell<TreeNode>>> = VecDeque::new();
    q.push_back(root.unwrap());
    while !q.is_empty() {
        let n = q.len();
        for i in 0..n {
            let node = q.pop_front().unwrap();
            let nb = node.borrow();
            if i == n - 1 {
                res.push(nb.val as i64);
            }
            if let Some(l) = nb.left.clone() {
                q.push_back(l);
            }
            if let Some(r) = nb.right.clone() {
                q.push_back(r);
            }
        }
    }
    res
}
