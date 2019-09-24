import React from "react";
import styled from "@emotion/styled/macro";

export const Alert = styled.p`
  color: var(--danger-color);
  font-weight: 700;
`;

export const Spinner = styled.div`
  display: inline-block;
  width: 4rem;
  height: 4rem;
  vertical-align: text-bottom;
  border: 0.25em solid white;
  border-right-color: white;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spinner-border 0.75s linear infinite;

  @keyframes spinner-border {
    to {
      transform: rotate(360deg);
    }
  }
`;

const ToolbarContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: baseline;
  justify-content: space-between;

  > div {
    position: relative;
    display: flex;
    margin-bottom: 2rem;

    > a {
      display: block;
      margin-left: 2rem;
      height: 24px;
      color: var(--secondary-color);
    }

    > a:first-of-type {
      margin-left: 0;
    }
  }

  @media (min-width: 700px) {
    flex-direction: row;

    > div {
      margin-bottom: 0;
    }
  }
`;

export const Toolbar = ({ title, children }) => {
  return (
    <ToolbarContainer>
      <h2>{title}</h2>
      <div>{children}</div>
    </ToolbarContainer>
  );
};
