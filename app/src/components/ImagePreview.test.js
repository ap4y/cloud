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
      albumName="bar"
      images={images}
      match={{
        path: "/test/:imageName",
        params: { imageName: image.name }
      }}
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
  expect(
    wrapper
      .find("AlbumItem")
      .first()
      .prop("src")
  ).toEqual("/api/gallery/bar/thumbnail/test/Test.jpg?jwt=foo");
});

it("renders image preview", () => {
  const wrapper = shallow(
    <ImagePreview
      images={images}
      match={{
        path: "/test/:imageName",
        params: { imageName: image.name }
      }}
      albumName="bar"
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
      match={{
        path: "/test/:imageName",
        params: { imageName: image.name }
      }}
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

it("renders exif data", async () => {
  const wrapper = shallow(
    <ImagePreview
      images={images}
      match={{
        path: "/test/:imageName",
        params: { imageName: image.name }
      }}
      fetchExif={() => new Promise(r => r({ foo: "bar" }))}
    />
  );

  wrapper
    .find("a")
    .findWhere(n => n.prop("href") === "#exif")
    .simulate("click", { preventDefault: () => {} });
  await new Promise(r => setTimeout(r, 100));

  expect(wrapper.exists("EXIFData")).toBeTruthy();
  expect(wrapper.find("EXIFData").prop("exif")).toEqual({ foo: "bar" });
});

it("renders thumbnail", () => {
  const wrapper = shallow(<AlbumItem image={image} src="test.jpg" />);

  wrapper.find("VisibilitySensor").invoke("onChange")(true);
  expect(wrapper.find("img").prop("src")).toEqual("test.jpg");
  expect(wrapper.find("figcaption").text()).toEqual("Test");
});
