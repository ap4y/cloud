import { configure } from "enzyme";
import Adapter from "enzyme-adapter-react-16";

import { apiClient } from "./lib/actions";

apiClient.authToken = "foo";

configure({ adapter: new Adapter() });
