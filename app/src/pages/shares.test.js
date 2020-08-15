import React from "react";
import { shallow } from "enzyme";
import { SharesList } from "./shares";
import GalleryItems from "../components/shares/GalleryItems";
import FilesItems from "../components/shares/FilesItems";

const share = {
  slug: "foo",
  name: "Test",
  expires_at: new Date(0),
  items: ["Test1.jpg", "Test2.jpg"],
  type: "gallery"
};

it("renders list of shares", () => {
  const wrapper = shallow(
    <SharesList shares={[share]} history={{ goBack: () => {} }} />
  );

  expect(wrapper.find("h3").length).toEqual(1);
  expect(
    wrapper
      .find("h3")
      .at(0)
      .text()
  ).toEqual("linkGallery: Test");

  expect(wrapper.find("GalleryItems").length).toEqual(1);

  expect(
    wrapper
      .find("div")
      .first()
      .find("div")
      .last()
      .text()
  ).toEqual(`access_time${new Date(0).toLocaleDateString()}`);
});

it("removes shares", () => {
  let removedSlug = null;
  const wrapper = shallow(
    <SharesList
      shares={[share]}
      removeShare={slug => (removedSlug = slug)}
      history={{ goBack: () => {} }}
    />
  );

  wrapper
    .find("button")
    .first()
    .simulate("click");
  expect(removedSlug).toEqual("foo");
});

it("renders gallery share items", () => {
  const wrapper = shallow(<GalleryItems gallery="foo" items={share.items} />);
  expect(wrapper.find("img").length).toEqual(2);
  expect(
    wrapper
      .find("img")
      .first()
      .prop("src")
  ).toEqual("/api/gallery/foo/thumbnail/Test1.jpg?jwt=foo");
});

it("renders files share items", () => {
  const wrapper = shallow(<FilesItems items={share.items} />);
  expect(wrapper.find("li").length).toEqual(2);
  expect(
    wrapper
      .find("li")
      .first()
      .text()
  ).toEqual("arrow_rightTest1.jpg");
});
