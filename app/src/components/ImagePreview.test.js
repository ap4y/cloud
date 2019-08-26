import React from "react";
import { shallow } from "enzyme";
import ImagePreview, { AlbumItem } from "./ImagePreview";

const image = {
  name: "Test",
  path: "test/Test.jpg",
  updated_at: new Date().toString()
};
const images = [
  image,
  {
    name: "Test 2",
    path: "test/Test 2.jpg",
    updated_at: new Date().toString()
  }
];

it("renders image thumbnails", () => {
  const wrapper = shallow(
    <ImagePreview
      images={images}
      match={{ path: `/test/:imageName`, params: { imageName: image.name } }}
      authToken="foo"
    />
  );

  expect(wrapper.find("AlbumItem").length).toEqual(2);
  expect(
    wrapper
      .find("AlbumItem")
      .first()
      .prop("image")
  ).toEqual(image);
  expect(
    wrapper
      .find("AlbumItem")
      .first()
      .prop("selected")
  ).toEqual(true);
});

it("renders image preview", () => {
  const wrapper = shallow(
    <ImagePreview
      images={images}
      match={{ path: `/test/:imageName`, params: { imageName: image.name } }}
      albumName="bar"
      authToken="foo"
    />
  );

  expect(
    wrapper
      .find("img")
      .first()
      .prop("src")
  ).toEqual("/api/gallery/bar/image/test/Test.jpg?jwt=foo");
});

it("renders navigation chevrons", () => {
  const wrapper = shallow(
    <ImagePreview
      images={images}
      match={{ path: `/test/:imageName`, params: { imageName: image.name } }}
      authToken="foo"
    />
  );

  expect(
    wrapper
      .find("Link")
      .first()
      .prop("to")
  ).toEqual("/test/Test 2");
  expect(
    wrapper
      .find("Link")
      .last()
      .prop("to")
  ).toEqual("/test/Test 2");
});

it("renders exif data", () => {
  const wrapper = shallow(
    <ImagePreview
      images={images}
      match={{ path: `/test/:imageName`, params: { imageName: image.name } }}
      authToken="foo"
    />
  );

  wrapper.setState({ exif: { foo: "bar" } });
  expect(wrapper.exists("EXIFData")).toBeTruthy();
  expect(wrapper.find("EXIFData").prop("exif")).toEqual({ foo: "bar" });
});

it("renders thumbnail", () => {
  const wrapper = shallow(
    <AlbumItem image={image} album="foo" authToken="foo" />
  );

  expect(wrapper.find("img").prop("src")).toEqual(
    "/api/gallery/foo/thumbnail/test/Test.jpg?jwt=foo"
  );
  expect(wrapper.find("figcaption").text()).toEqual("Test");
});
