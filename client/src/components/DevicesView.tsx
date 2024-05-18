import axios from "axios";
import { useEffect, useState } from "react";

interface Device {
  Imei: string;
  BatteryPower: number;
  Charging: boolean;
}

function DevicesView() {
  const [devices, setDevices] = useState<[Device]>();

  useEffect(() => {
    axios
      .get("http://localhost:50731/v1/devices", { withCredentials: true })
      .then((response) => {
        if (response.status === 200) {
          setDevices(response.data);
        }
        //todo: else toast
      })
      .catch((err) => {
        console.error(err);
        //todo: toast
      });
  }, []);

  return (
    <div className="flex flex-col gap-20 w-full">
      <p className="text-5xl font-semibold">Your Devices:</p>
      <div className="px-20 flex flex-row gap-20 w-full">
        {devices?.map((device) => (
          <div key={device.Imei}>{device.Imei}</div>
        ))}
      </div>
      <div className="flex flex-col items-start gap-2">
        <p className="text-gray-800 text-lg">Can not find your device?</p>
        <button className="bg-gray-300 px-10 py-3 rounded-md shadow-md font-bold hover:scale-110 duration-100">
          Link your device
        </button>
      </div>
    </div>
  );
  // todo: link device
}

export default DevicesView;
