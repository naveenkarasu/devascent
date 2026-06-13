use std::rc::Rc;
use std::cell::RefCell;
use std::collections::VecDeque;

fn level_order(root: Option<Rc<RefCell<TreeNode>>>) -> Vec<Vec<i64>> {
    let mut res: Vec<Vec<i64>> = Vec::new();
    if root.is_none() {
        return res;
    }
    let mut q: VecDeque<Rc<RefCell<TreeNode>>> = VecDeque::new();
    q.push_back(root.unwrap());
    while !q.is_empty() {
        let n = q.len();
        let mut level: Vec<i64> = Vec::new();
        for _ in 0..n {
            let node = q.pop_front().unwrap();
            let nb = node.borrow();
            level.push(nb.val as i64);
            if let Some(l) = nb.left.clone() {
                q.push_back(l);
            }
            if let Some(r) = nb.right.clone() {
                q.push_back(r);
            }
        }
        res.push(level);
    }
    res
}
