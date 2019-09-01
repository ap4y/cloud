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
      albumName="Test"
      authToken="foo"
      images={images}
      match={{ url: "/test", params: { albumName: "bar" } }}
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
  expect(
    wrapper
      .find("ImageCell")
      .first()
      .prop("src")
  ).toEqual("/api/test/thumbnail/test/Test.jpg?jwt=foo");
});

it("renders image cell", () => {
  const wrapper = shallow(<ImageCell image={image} src="test.jpg" />);

  expect(wrapper.find("figcaption").text()).toEqual(
    "Test1/1/1970, 12:00:00 PM"
  );
  expect(wrapper.find("img").prop("src")).toEqual("test.jpg");
});

it("toggles share popup", () => {
  const images = [image];
  const wrapper = shallow(
    <ImageGrid
      albumName="Test"
      images={images}
      match={{ url: "/test", params: { albumName: "bar" } }}
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
      albumName="Test"
      images={images}
      match={{ url: "/test", params: { albumName: "bar" } }}
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

it("requests shares", () => {
  let result = {};
  const wrapper = shallow(
    <ImageGrid
      albumName="Test"
      share="foo"
      images={[]}
      match={{ url: "/test", params: { albumName: "bar" } }}
      fetchAlbum={(album, share) => {
        result = { album, share };
      }}
    />
  );

  expect(result.album).toEqual("Test");
  expect(result.share).toEqual("foo");
});
