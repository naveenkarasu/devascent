use std::rc::Rc;
use std::cell::RefCell;

fn diameter_of_binary_tree(root: Option<Rc<RefCell<TreeNode>>>) -> i64 {
    fn height(node: &Option<Rc<RefCell<TreeNode>>>, best: &mut i64) -> i64 {
        match node {
            None => 0,
            Some(n) => {
                let nb = n.borrow();
                let l = height(&nb.left, best);
                let r = height(&nb.right, best);
                if l + r > *best {
                    *best = l + r;
                }
                1 + l.max(r)
            }
        }
    }
    let mut best = 0i64;
    height(&root, &mut best);
    best
}
