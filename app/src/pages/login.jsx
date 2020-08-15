import React, { useState } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";

import { Alert } from "../components/Controls";

import { signIn, resetAuthError, resetErrorMessage } from "../lib/actions";

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;

  form {
    width: 100%;
    max-width: 30rem;
  }
`;

export const LoginForm = ({
  authSuccess,
  errorMessage,
  location,
  signIn,
  resetErrorMessage,
  resetAuthError
}) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const { from } = location.state || { from: { pathname: "/" } };
  if (authSuccess) return <Redirect to={from} />;

  const performAuth = e => {
    e.preventDefault();
    resetAuthError();
    resetErrorMessage();
    signIn(username, password);
  };

  return (
    <Container>
      <form onSubmit={performAuth}>
        <label>Username: </label>
        <input
          type="text"
          onChange={({ target }) => setUsername(target.value)}
        />
        <label>Password:</label>
        <input
          type="password"
          onChange={({ target }) => setPassword(target.value)}
        />
        {errorMessage && <Alert>Invalid username or password.</Alert>}
        <input type="submit" value="Login" />
      </form>
    </Container>
  );
};

export default connect(
  ({ auth, errorMessage }) => ({ authSuccess: auth.success, errorMessage }),
  { signIn, resetErrorMessage, resetAuthError }
)(LoginForm);
