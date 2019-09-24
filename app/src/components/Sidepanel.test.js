import React from "react";
import { shallow } from "enzyme";
import Sidepanel from "./Sidepanel";

it("renders sidepanel", () => {
  let signedOut = false;
  const wrapper = shallow(
    <Sidepanel
      navItems={[{ title: "Test", route: "/test", icon: "folder" }]}
      onSignOut={() => (signedOut = true)}
    >
      <div>Test</div>
    </Sidepanel>
  );

  expect(wrapper.find("nav").text()).toEqual("Test");
  expect(wrapper.find("NavLink").length).toEqual(2);

  expect(
    wrapper
      .find("NavLink")
      .at(0)
      .prop("to")
  ).toEqual("/test");
  expect(
    wrapper
      .find("NavLink")
      .at(1)
      .prop("to")
  ).toEqual("/shares");

  wrapper.find("a").simulate("click");
  expect(signedOut).toEqual(true);
});
