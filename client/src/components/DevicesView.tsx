import axios from "axios";
import { useEffect, useState } from "react";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { toast } from "sonner";

interface Device {
  Imei: string;
  BatteryPower: number;
  Charging: boolean;
}

function DevicesView() {
  const [devices, setDevices] = useState<[Device]>();
  const [imei, setImei] = useState("");

  useEffect(() => {
    axios
      .get("http://localhost:50731/api/v1/devices", { withCredentials: true })
      .then((response) => {
        setDevices(response.data);
      })
      .catch(() => {
        toast("Something went wrong while quering devices");
      });
  }, []);

  const handle_link = () => {
    axios
      .get("http://localhost:50731/api/v1/link/" + imei, {
        withCredentials: true,
      })
      .then(() => {
        toast("Linked, refresh this page");
      })
      .catch(() => {
        toast("Something went wrong, check provided IMEI");
      });
  };

  return (
    <div className="flex flex-col gap-20 w-full">
      <p className="text-3xl font-semibold">Your Devices:</p>
      <div className="px-20 flex flex-row gap-20 w-full">
        {devices?.map((device) => (
          <div key={device.Imei}>{device.Imei}</div>
        ))}
      </div>
      <div className="flex flex-col items-start gap-2">
        <Popover>
          <PopoverTrigger asChild>
            <p className="text-blue-400 text-lg underline cursor-pointer">
              Can not find your device?
            </p>
          </PopoverTrigger>
          <PopoverContent className="w-96 ml-2">
            <div className="grid gap-4">
              <div className="space-y-2">
                <h4 className="font-medium leading-none">Link It!</h4>
              </div>
              <div className="grid gap-2">
                <div className="grid grid-cols-3 items-center gap-4">
                  <Label>IMEI</Label>
                  <Input
                    value={imei}
                    onChange={(e) => setImei(e.target.value)}
                    className="col-span-2 h-8"
                  />
                </div>
              </div>
              <Button
                className="dark border border-blue-100"
                onClick={handle_link}
              >
                Link
              </Button>
            </div>
          </PopoverContent>
        </Popover>
      </div>
    </div>
  );
}

export default DevicesView;
