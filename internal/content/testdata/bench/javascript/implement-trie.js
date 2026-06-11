function trie_ops(operations) {
  const root = {};

  function insert(word) {
    let node = root;
    for (const c of word) {
      if (!node[c]) node[c] = {};
      node = node[c];
    }
    node['$'] = true;
  }

  function search(word) {
    let node = root;
    for (const c of word) {
      if (!node[c]) return false;
      node = node[c];
    }
    return node['$'] === true;
  }

  function startsWith(prefix) {
    let node = root;
    for (const c of prefix) {
      if (!node[c]) return false;
      node = node[c];
    }
    return true;
  }

  const out = [];
  for (const op of operations) {
    const name = op[0];
    const arg = op[1];
    if (name === 'insert') {
      insert(arg);
      out.push(null);
    } else if (name === 'search') {
      out.push(search(arg));
    } else {
      out.push(startsWith(arg));
    }
  }
  return out;
}
