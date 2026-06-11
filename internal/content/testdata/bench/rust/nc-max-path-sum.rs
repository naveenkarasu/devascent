use std::rc::Rc;
use std::cell::RefCell;

fn gain(node: &Option<Rc<RefCell<TreeNode>>>, best: &mut i64) -> i64 {
    match node {
        None => 0,
        Some(n) => {
            let nb = n.borrow();
            let l = gain(&nb.left, best).max(0);
            let r = gain(&nb.right, best).max(0);
            let val = nb.val as i64;
            if val + l + r > *best {
                *best = val + l + r;
            }
            val + l.max(r)
        }
    }
}

fn max_path_sum(root: Option<Rc<RefCell<TreeNode>>>) -> i64 {
    let mut best: i64 = -1_000_000_000;
    gain(&root, &mut best);
    best
}
