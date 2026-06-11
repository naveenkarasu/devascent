use std::rc::Rc;
use std::cell::RefCell;

fn lca_bst(root: Option<Rc<RefCell<TreeNode>>>, p: i64, q: i64) -> i64 {
    let mut node = root;
    while let Some(n) = node {
        let v = n.borrow().val as i64;
        if p < v && q < v {
            let left = n.borrow().left.clone();
            node = left;
        } else if p > v && q > v {
            let right = n.borrow().right.clone();
            node = right;
        } else {
            return v;
        }
    }
    -1
}
