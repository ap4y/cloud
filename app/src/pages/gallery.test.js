import React from "react";
import { shallow } from "enzyme";
import { ImageGrid, ImageCell } from "./gallery";

const image = {
  name: "Test",
  path: "test/Test.jpg",
  updated_at: new Date(0).toString()
};

it("renders image list", () => {
  const images = [
    image,
    {
      name: "Test 2",
      path: "test/Test 2.jpg",
      updated_at: new Date().toString()
    }
  ];
  const wrapper = shallow(
    <ImageGrid
      galleryName="Test"
      images={images}
      match={{ url: "/test", params: { galleryName: "bar" } }}
      fetchAlbum={() => {}}
    />
  );

  expect(wrapper.find("h2").text()).toEqual("Test");
  expect(wrapper.find("ImageCell").length).toEqual(2);
  expect(
    wrapper
      .find("ImageCell")
      .first()
      .prop("image")
  ).toEqual(image);
});

it("renders image cell", () => {
  const wrapper = shallow(
    <ImageCell image={image} gallery="bar" authToken="foo" />
  );

  expect(wrapper.find("figcaption").text()).toEqual(
    "Test1/1/1970, 12:00:00 PM"
  );
  expect(wrapper.find("img").prop("src")).toEqual(
    "/api/gallery/bar/thumbnail/test/Test.jpg?jwt=foo"
  );
});

it("toggles share popup", () => {
  const images = [image];
  const wrapper = shallow(
    <ImageGrid
      galleryName="Test"
      images={images}
      match={{ url: "/test", params: { galleryName: "bar" } }}
      fetchAlbum={() => {}}
    />
  );

  expect(wrapper.find("SharePopup").exists()).toBeFalsy();
  wrapper.setState({ showSharing: true });
  expect(wrapper.find("SharePopup").exists()).toBeTruthy();
});

it("creates shares", () => {
  let shared = false;
  const images = [image];
  const wrapper = shallow(
    <ImageGrid
      galleryName="Test"
      images={images}
      match={{ url: "/test", params: { galleryName: "bar" } }}
      fetchAlbum={() => {}}
      shareAlbum={() => {
        shared = true;
        return new Promise(() => {});
      }}
    />
  );

  wrapper.setState({ showSharing: true });
  expect(wrapper.find("SharePopup").exists()).toBeTruthy();
  wrapper.find("SharePopup").invoke("onShare")();
  expect(shared).toBeTruthy();
});
