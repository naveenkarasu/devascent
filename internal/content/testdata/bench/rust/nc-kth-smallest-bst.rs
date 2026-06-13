use std::rc::Rc;
use std::cell::RefCell;

fn kth_smallest(root: Option<Rc<RefCell<TreeNode>>>, k: i64) -> i64 {
    let mut stack: Vec<Rc<RefCell<TreeNode>>> = Vec::new();
    let mut cur = root;
    let mut k = k;
    while !stack.is_empty() || cur.is_some() {
        while let Some(n) = cur {
            let left = n.borrow().left.clone();
            stack.push(n);
            cur = left;
        }
        let node = stack.pop().unwrap();
        k -= 1;
        if k == 0 {
            return node.borrow().val as i64;
        }
        cur = node.borrow().right.clone();
    }
    -1
}
