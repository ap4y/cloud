import React from "react";
import { shallow } from "enzyme";
import SharePopup from "./SharePopup";

it("renders popup", () => {
  const wrapper = shallow(<SharePopup items={["foo", "bar"]} />);

  expect(wrapper.find("h4").text()).toEqual("Sharing 2 items");
  expect(wrapper.find("button").exists()).toBeTruthy();
});

it("renders errors", () => {
  const wrapper = shallow(
    <SharePopup items={["foo", "bar"]} error="Test error" />
  );

  expect(wrapper.childAt(2).text()).toEqual("Test error");
  expect(wrapper.find("button").exists()).toBeTruthy();
});

it("renders share link", () => {
  const wrapper = shallow(<SharePopup items={["foo", "bar"]} slug="foo" />);

  expect(wrapper.find("p").text()).toEqual("Share Link");
  expect(wrapper.find("a").prop("href")).toEqual("http://localhost/share/foo");
  expect(wrapper.find("button").exists()).toBeFalsy();
});

it("shares link", () => {
  let res = null;
  const wrapper = shallow(
    <SharePopup
      items={["foo", "bar"]}
      onShare={expiresAt => (res = expiresAt)}
    />
  );

  wrapper.find("button").simulate("click");
  expect(res).not.toBeNull();
});
