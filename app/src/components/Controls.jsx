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
