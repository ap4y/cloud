import { locateInTree, removeFromTree, addToTree } from "./utils";

const file = { name: "foo", path: "/foo", type: "file" };
const subfolderFile = { name: "baz", path: "/bar/baz", type: "file" };
const subfolder = {
  name: "bar",
  path: "/bar",
  type: "directory",
  children: [subfolderFile]
};
const tree = {
  name: "/",
  path: "/",
  children: [file, subfolder]
};

describe("locateInTree", () => {
  it("returns root", () => {
    const { folder, file } = locateInTree(tree, "/");
    expect(folder).toEqual(tree);
    expect(file).toEqual(null);
  });

  it("returns file in root", () => {
    const { folder, file } = locateInTree(tree, "/foo");
    expect(folder).toEqual(tree);
    expect(file).toEqual(file);
  });

  it("returns subfolder", () => {
    const { folder, file } = locateInTree(tree, "/bar");
    expect(folder).toEqual(subfolder);
    expect(file).toEqual(null);
  });

  it("returns file in subfolder", () => {
    const { folder, file } = locateInTree(tree, "/bar/baz");
    expect(folder).toEqual(subfolder);
    expect(file).toEqual(subfolderFile);
  });
});

describe("removeFromTree", () => {
  it("removes root file", () => {
    const newTree = removeFromTree(tree, { path: "/" }, file);
    expect(newTree.children.length).toEqual(1);
    expect(tree.children.length).toEqual(2);
  });

  it("removes subfolder file", () => {
    const newTree = removeFromTree(tree, { path: "/bar" }, subfolderFile);
    expect(newTree.children.length).toEqual(2);
    expect(newTree.children[1].children.length).toEqual(0);
  });
});

describe("addToTree", () => {
  it("adds root file", () => {
    const newTree = addToTree(tree, { path: "/" }, { name: "test" });
    expect(newTree.children.length).toEqual(3);
    expect(tree.children.length).toEqual(2);
  });

  it("adds subfolder file", () => {
    const newTree = addToTree(tree, { path: "/bar" }, { name: "test" });
    expect(newTree.children.length).toEqual(2);
    expect(newTree.children[1].children.length).toEqual(1);
  });
});
