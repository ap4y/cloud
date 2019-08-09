import React from "react";
import { shallow } from "enzyme";
import ImageEXIF from "./ImageEXIF";

it("renders exif data", () => {
  const exif = {
    DateTime: "2019/01/01",
    PixelXDimension: 20,
    PixelYDimension: 30,
    Make: "Sony",
    Model: "a6000",
    LensModel: "Sigma 30mm",
    ExposureTime: "1/30",
    FNumber: "1/100",
    ISOSpeedRatings: "100",
    MeteringMode: 5,
    Flash: 0x4f
  };
  const wrapper = shallow(<ImageEXIF exif={exif} />);

  expect(
    wrapper
      .find("li")
      .at(0)
      .text()
  ).toEqual("Created: 2019/01/01");

  expect(
    wrapper
      .find("li")
      .at(1)
      .text()
  ).toEqual("Dimensions: 20x30");

  expect(
    wrapper
      .find("li")
      .at(2)
      .text()
  ).toEqual("Camera: Sony a6000");

  expect(
    wrapper
      .find("li")
      .at(3)
      .text()
  ).toEqual("Lens: Sigma 30mm");

  expect(
    wrapper
      .find("li")
      .at(4)
      .text()
  ).toEqual("Exposure: 1/30");

  expect(
    wrapper
      .find("li")
      .at(5)
      .text()
  ).toEqual("Focal Length: 1/100");

  expect(
    wrapper
      .find("li")
      .at(6)
      .text()
  ).toEqual("ISO: 100");

  expect(
    wrapper
      .find("li")
      .at(7)
      .text()
  ).toEqual("Metering: Pattern");

  expect(
    wrapper
      .find("li")
      .at(8)
      .text()
  ).toEqual(
    "Flash: Flash fired, compulsory flash mode, red-eye reduction mode, return light detected"
  );
});
