export function locateInTree(tree, path) {
  const components = path.split("/");

  let folder = tree;
  let file = null;

  components.forEach(component => {
    const item = folder.children.find(({ name }) => name === component);
    if (!item) return;

    if (item.type === "file") {
      file = item;
    } else {
      folder = item;
    }
  });

  return { folder, file };
}

export function removeFromTree(tree, folder, file) {
  if (folder.path === "/") {
    return {
      ...tree,
      children: tree.children.filter(({ name }) => name !== file.name)
    };
  }

  const components = folder.path.split("/").filter(c => c.length > 0);
  let curFolder = tree;
  components.forEach((component, idx) => {
    const item = curFolder.children.find(({ name }) => name === component);
    if (!item) return;

    if (idx === components.length - 1) {
      item.children = item.children.filter(({ name }) => name !== file.name);
    }

    curFolder = item;
  });

  return tree;
}

export function addToTree(tree, folder, file) {
  removeFromTree(tree, folder, file);

  if (folder.path === "/") {
    return { ...tree, children: [...tree.children, file] };
  }

  const components = folder.path.split("/").filter(c => c.length > 0);
  let curFolder = tree;
  components.forEach((component, idx) => {
    const item = curFolder.children.find(({ name }) => name === component);
    if (!item) return;

    if (idx === components.length - 1) {
      item.children = [...item.children, file];
    }

    curFolder = item;
  });

  return tree;
}
