import React from "react";
import { shallow } from "enzyme";
import { LoginForm } from "./login";

it("renders form", () => {
  shallow(<LoginForm location={{}} />);
});

it("redirect to from with authToken", () => {
  const wrapper = shallow(
    <LoginForm location={{ state: { from: "/foo" } }} authToken="bar" />
  );
  expect(wrapper.exists("Redirect")).toBeTruthy();
  expect(wrapper.find("Redirect").prop("to")).toEqual("/foo");
});

it("renders error message", () => {
  const wrapper = shallow(
    <LoginForm location={{}} errorMessage="Invalid login" />
  );
  expect(
    wrapper
      .find("form")
      .children()
      .at(4)
      .text()
  ).toEqual("Invalid username or password.");
});

it("performs auth", () => {
  let state = {};

  const wrapper = shallow(
    <LoginForm
      location={{}}
      resetErrorMessage={() => {
        state.resetError = true;
      }}
      resetAuthError={() => {
        state.resetAuth = true;
      }}
      signIn={(username, password) => {
        state.username = username;
        state.password = password;
      }}
    />
  );

  wrapper
    .find("form")
    .children()
    .at(1)
    .simulate("change", { target: { value: "foo" } });
  wrapper
    .find("form")
    .children()
    .at(3)
    .simulate("change", { target: { value: "bar" } });
  wrapper.find("form").simulate("submit");

  expect(state.username).toEqual("foo");
  expect(state.password).toEqual("bar");
  expect(state.resetError).toBeTruthy();
  expect(state.resetAuth).toBeTruthy();
});
