import React from "react";
import { act } from "react-dom/test-utils";
import { shallow, mount } from "enzyme";
import { FilesGrid, DirCell, FileCell } from "./files";

const tree = {
  name: "test",
  path: "/test",
  children: [
    { name: "foo", path: "/foo", type: "file" },
    { name: "bar", path: "/bar", type: "directory" }
  ]
};

it("renders files list", () => {
  const wrapper = shallow(<FilesGrid folder={tree} />);

  expect(wrapper.find("FilesToolbar").prop("path")).toEqual("/test");

  expect(wrapper.find("DirCell").length).toEqual(1);
  expect(wrapper.find("DirCell").prop("dir")).toEqual({
    name: "bar",
    path: "/bar",
    type: "directory"
  });

  expect(wrapper.find("FileCell").length).toEqual(1);
  expect(wrapper.find("FileCell").prop("file")).toEqual({
    name: "foo",
    path: "/foo",
    type: "file"
  });

  expect(
    wrapper
      .find("NavLink")
      .first()
      .prop("to")
  ).toEqual("/files/bar/");
  expect(
    wrapper
      .find("NavLink")
      .last()
      .prop("to")
  ).toEqual("/files/foo");
});

it("renders specific file", () => {
  const file = { name: "foo", path: "/foo" };
  const wrapper = shallow(<FilesGrid folder={tree} file={file} />);

  expect(wrapper.find("FilesToolbar").prop("path")).toEqual("/test");
  expect(wrapper.find("FilesToolbar").prop("file")).toEqual(file);

  expect(wrapper.find("DirCell").length).toEqual(0);
  expect(wrapper.find("FileCell").length).toEqual(1);

  expect(wrapper.find("FileCell").prop("file")).toEqual(file);
  expect(wrapper.find("FileCell").prop("large")).toEqual(true);
  expect(wrapper.find("FileCell").prop("withModTime")).toEqual(true);
});

it("renders viewable file content", async () => {
  const file = {
    name: "foo.org",
    path: "/foo.org",
    url: "/files/test/files/foo.org"
  };

  let fetchedURL = null;
  let wrapper = null;
  await act(async () => {
    wrapper = mount(
      <FilesGrid
        folder={tree}
        file={file}
        fetchFile={url => {
          fetchedURL = url;
          return new Promise(resolve => resolve("foo"));
        }}
      />
    );
  });

  wrapper.update();

  expect(fetchedURL).toEqual(file.url);

  expect(wrapper.find("FilesToolbar").prop("path")).toEqual("/test");
  expect(wrapper.find("FilesToolbar").prop("file")).toEqual(file);

  expect(wrapper.find("TextEditor").length).toEqual(1);
  expect(wrapper.find("TextEditor").prop("value")).toEqual("foo");
});

it("renders dir cell", () => {
  const wrapper = shallow(
    <DirCell
      dir={{
        name: "foo",
        path: "/foo"
      }}
    />
  );

  expect(wrapper.find("span").text()).toEqual("foo");
  expect(wrapper.find("i").text()).toEqual("folder");
});

it("renders file cell", () => {
  const wrapper = shallow(
    <FileCell
      file={{
        name: "foo",
        path: "/foo",
        updated_at: new Date(0).toString()
      }}
      withModTime
    />
  );

  expect(
    wrapper
      .find("span")
      .first()
      .text()
  ).toEqual("foo");
  expect(wrapper.find("i").text()).toEqual("insert_drive_file");
  expect(
    wrapper
      .find("span")
      .last()
      .text()
  ).toEqual("Updated At: 1/1/1970, 12:00:00 PM");
});

it("uploads file", () => {
  let result = {};
  const wrapper = shallow(
    <FilesGrid
      folder={tree}
      uploadFile={(folder, file) => {
        result.folder = folder;
        result.file = file;
      }}
    />
  );
  wrapper.find("input").simulate("change", { target: { files: ["foo"] } });

  expect(result.folder).toEqual(tree);
  expect(result.file).toEqual("foo");
});

it("saves file", () => {
  let result = {};
  const wrapper = shallow(
    <FilesGrid
      folder={tree}
      file={{ name: "foo" }}
      uploadFile={(folder, file) => {
        result.folder = folder;
        result.file = file;
      }}
    />
  );
  wrapper.find("FilesToolbar").invoke("onSave")({ preventDefault: () => {} });

  expect(result.folder).toEqual(tree);
  expect(result.file.name).toEqual("foo");
});

it("removes file", () => {
  global.confirm = jest.fn(() => true);

  let result = {};
  const wrapper = shallow(
    <FilesGrid
      folder={tree}
      file={{ name: "foo" }}
      removeFile={(folder, file) => {
        result.folder = folder;
        result.file = file;
        return new Promise(() => {});
      }}
    />
  );
  wrapper.find("FilesToolbar").invoke("onDelete")({ preventDefault: () => {} });

  expect(result.folder).toEqual(tree);
  expect(result.file.name).toEqual("foo");
});
