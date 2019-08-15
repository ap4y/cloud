import React from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";

import { Alert } from "../components/Typography";

import { signIn, resetAuthError, resetErrorMessage } from "../actions";

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

export class LoginForm extends React.Component {
  state = { username: "", password: "" };

  _handleChange = field => ({ target }) => {
    this.setState({ [field]: target.value });
  };

  performAuth = () => {
    const { username, password } = this.state;
    this.props.resetAuthError();
    this.props.resetErrorMessage();
    this.props.signIn(username, password);
  };

  render() {
    const { from } = this.props.location.state || { from: { pathname: "/" } };
    const { errorMessage } = this.props;
    if (this.props.authToken) return <Redirect to={from} />;

    return (
      <Container>
        <form onSubmit={this.performAuth}>
          <label>Username: </label>
          <input type="text" onChange={this._handleChange("username")} />
          <label>Password:</label>
          <input type="password" onChange={this._handleChange("password")} />
          {errorMessage && <Alert>Invalid username or password.</Alert>}
          <input type="submit" value="Login" />
        </form>
      </Container>
    );
  }
}

export default connect(
  ({ authToken, errorMessage }) => ({ authToken, errorMessage }),
  { signIn, resetErrorMessage, resetAuthError }
)(LoginForm);