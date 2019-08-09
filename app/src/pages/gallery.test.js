import React from "react";
import { shallow } from "enzyme";
import { ImageGrid, ImageCell } from "./gallery";

it("renders image list", () => {
  const image = {
    name: "Test",
    path: "/test/Test.jpg",
    updated_at: new Date().toString()
  };
  const images = [
    image,
    {
      name: "Test 2",
      path: "/test/Test 2.jpg",
      updated_at: new Date().toString()
    }
  ];
  const wrapper = shallow(
    <ImageGrid
      galleryName="Test"
      images={images}
      match={{ url: "/test" }}
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
  const image = {
    name: "Test",
    path: "/test/Test.jpg",
    updated_at: new Date(0).toString()
  };

  const wrapper = shallow(<ImageCell image={image} authToken="foo" />);

  expect(wrapper.find("figcaption").text()).toEqual(
    "Test1/1/1970, 12:00:00 PM"
  );
  expect(wrapper.find("img").prop("src")).toEqual(
    "/api/gallery/thumbnails//test/Test.jpg?jwt=foo"
  );
});
