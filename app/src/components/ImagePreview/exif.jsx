import React from "react";
import styled from "@emotion/styled";

const meteringModes = {
  0: "Unknown",
  1: "Average",
  2: "CenterWeightedAverage",
  3: "Spot",
  4: "MultiSpot",
  5: "Pattern",
  6: "Partial",
  255: "other"
};

const flashModes = {
  0x0: "Flash did not fire",
  0x1: "Flash fired",
  0x5: "Strobe return light not detected",
  0x7: "Strobe return light detected",
  0x9: "Flash fired, compulsory flash mode",
  0xd: "Flash fired, compulsory flash mode, return light not detected",
  0xf: "Flash fired, compulsory flash mode, return light detected",
  0x10: "Flash did not fire, compulsory flash mode",
  0x18: "Flash did not fire, auto mode",
  0x19: "Flash fired, auto mode",
  0x1d: "Flash fired, auto mode, return light not detected",
  0x1f: "Flash fired, auto mode, return light detected",
  0x20: "No flash function",
  0x41: "Flash fired, red-eye reduction mode",
  0x45: "Flash fired, red-eye reduction mode, return light not detected",
  0x47: "Flash fired, red-eye reduction mode, return light detected",
  0x49: "Flash fired, compulsory flash mode, red-eye reduction mode",
  0x4d: "Flash fired, compulsory flash mode, red-eye reduction mode, return light not detected",
  0x4f: "Flash fired, compulsory flash mode, red-eye reduction mode, return light detected",
  0x59: "Flash fired, auto mode, red-eye reduction mode",
  0x5d: "Flash fired, auto mode, return light not detected, red-eye reduction mode",
  0x5f: "Flash fired, auto mode, return light detected, red-eye reduction mode"
};

const EXIFList = styled.ul`
  list-style: none;
  padding: 2rem;
  margin: 0;
  background: #3b4252;

  color: #eceff4;
`;

const EXIFData = ({ exif }) => {
  return (
    <EXIFList>
      <li>
        <strong>Created:</strong> {exif.DateTime}
      </li>
      <li>
        <strong>Dimensions:</strong> {exif.PixelXDimension}x
        {exif.PixelYDimension}
      </li>
      <li>
        <strong>Camera:</strong> {exif.Make} {exif.Model}
      </li>
      <li>
        <strong>Lens:</strong> {exif.LensModel}
      </li>
      <li>
        <strong>Exposure:</strong> {exif.ExposureTime}
      </li>
      <li>
        <strong>Focal Length:</strong> {exif.FNumber}
      </li>
      <li>
        <strong>ISO:</strong> {exif.ISOSpeedRatings}
      </li>
      <li>
        <strong>Metering:</strong> {meteringModes[exif.MeteringMode]}
      </li>
      <li>
        <strong>Flash:</strong> {flashModes[exif.Flash]}
      </li>
    </EXIFList>
  );
};

export default EXIFData;
