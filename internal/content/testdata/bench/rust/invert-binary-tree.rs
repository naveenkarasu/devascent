use std::rc::Rc;
use std::cell::RefCell;

fn invert_tree(root: Option<Rc<RefCell<TreeNode>>>) -> Option<Rc<RefCell<TreeNode>>> {
    if let Some(node) = &root {
        let left = node.borrow_mut().left.take();
        let right = node.borrow_mut().right.take();
        node.borrow_mut().left = invert_tree(right);
        node.borrow_mut().right = invert_tree(left);
    }
    root
}
