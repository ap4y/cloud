import React from "react";
import { shallow } from "enzyme";
import { AlbumsList } from "./albums";

it("renders list of albums", () => {
  const albums = [
    { name: "Test", items_count: 2 },
    { name: "Test 2", items_count: 3 }
  ];
  const wrapper = shallow(
    <AlbumsList albums={albums} match={{ url: "/test" }} />
  );

  expect(wrapper.find("li").length).toEqual(2);
  expect(
    wrapper
      .find("li NavLink")
      .at(0)
      .prop("to")
  ).toEqual("/test/Test");

  expect(
    wrapper
      .find("li NavLink")
      .at(0)
      .children()
      .at(0)
      .text()
  ).toEqual("Test");
  expect(
    wrapper
      .find("li NavLink")
      .at(0)
      .children()
      .at(1)
      .text()
  ).toEqual("2");
});
