import React, { Component } from "react";
import styled from "@emotion/styled";

import { Alert } from "./Typography";

const SharePopupContainer = styled.div`
  width: 30rem;
  padding: 2.5rem 2.5rem 1.5rem 2.5rem;
  z-index: 10;

  background: white;
  border: 1px solid var(--outline-color);
  border-radius: 5px;

  text-align: center;

  div {
    display: flex;
    align-items: baseline;
  }

  input {
    flex: 1;
    margin-left: 1rem;
  }

  &:after {
    content: "";
    position: absolute;
    top: -10px;
    right: 10px;
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 0 10px 10px 10px;
    border-color: transparent transparent var(--outline-color) transparent;
  }
`;

class SharePopup extends Component {
  state = { expireAt: new Date() };

  componentDidMount() {
    let expireAt = new Date();
    expireAt.setDate(expireAt.getDate() + 7);
    this.setState({ expireAt });
  }

  updateExpireAt = ({ target }) => {
    this.setState({ expireAt: new Date(target.value) });
  };

  render() {
    const { items, onShare, error, slug } = this.props;
    const { expireAt } = this.state;

    return (
      <SharePopupContainer>
        <h4>{`Sharing ${items.length} ${
          items.length === 1 ? "item" : "items"
        }`}</h4>
        <div>
          <label>Expire At:</label>
          <input
            type="date"
            value={expireAt.toISOString().substr(0, 10)}
            onChange={this.updateExpireAt}
          />
        </div>

        {slug && (
          <p>
            <a
              href={`${window.location.origin}/share/${slug}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Share Link
            </a>
          </p>
        )}
        {error && <Alert>{error}</Alert>}

        {!slug && <button onClick={() => onShare(expireAt)}>Share</button>}
      </SharePopupContainer>
    );
  }
}

export default SharePopup;
