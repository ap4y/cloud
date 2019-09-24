import React from "react";
import { shallow } from "enzyme";
import { FilesTree, DirTree } from "./files_tree";

const tree = {
  name: "test",
  path: "/test",
  children: [
    { name: "foo", path: "/foo", type: "file" },
    { name: "bar", path: "/bar", type: "directory" }
  ]
};

it("renders root tree item", () => {
  const wrapper = shallow(<FilesTree tree={tree} />);

  expect(wrapper.find("DirTree").prop("leaf")).toEqual(tree);
  expect(wrapper.find("DirTree").prop("level")).toEqual(0);
});

it("renders tree item and it's children", () => {
  const wrapper = shallow(<DirTree leaf={tree} level={0} />);

  expect(wrapper.find("DirItem").length).toEqual(1);
  expect(wrapper.find("DirItem NavLink").prop("to")).toEqual("/files/test/");
  expect(wrapper.find("DirItem span").text()).toEqual("test");

  expect(wrapper.find("DirTree").length).toEqual(1);
  expect(wrapper.find("DirTree").prop("leaf")).toEqual(tree.children[1]);
  expect(wrapper.find("DirTree").prop("level")).toEqual(1);
});
