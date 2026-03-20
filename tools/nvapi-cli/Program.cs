using System;
using System.IO;
using System.Linq;
using NvAPIWrapper;
using NvAPIWrapper.DRS;

namespace NvApiCli
{
    class Program
    {
        const uint NGX_DLSS_SR_OVERRIDE_RENDER_PRESET_SELECTION_ID = 0x10E41DF3;
        const uint NGX_DLSS_RR_OVERRIDE_RENDER_PRESET_SELECTION_ID = 0x10E41DF7;

        static int Main(string[] args)
        {
            if (args.Length < 2)
            {
                Console.WriteLine("{\"success\": false, \"error\": \"Usage: <get|set> <exeName> [tech] [preset]\"}");
                return 1;
            }

            string command = args[0].ToLowerInvariant();
            string exeName = args[1];

            try
            {
                NVIDIA.Initialize();
                var session = DriverSettingsSession.CreateAndLoad();
                var profile = FindProfile(session, exeName);

                if (command == "get")
                {
                    if (profile == null)
                    {
                        Console.WriteLine("{\"success\": true, \"dlssPreset\": 0, \"dlssdPreset\": 0, \"foundProfile\": false}");
                    }
                    else
                    {
                        uint dlssPreset = 0;
                        var srSetting = profile.Settings.FirstOrDefault(x => x.SettingId == NGX_DLSS_SR_OVERRIDE_RENDER_PRESET_SELECTION_ID);
                        if (srSetting != null && srSetting.CurrentValue is uint)
                            dlssPreset = (uint)srSetting.CurrentValue;

                        uint dlssdPreset = 0;
                        var rrSetting = profile.Settings.FirstOrDefault(x => x.SettingId == NGX_DLSS_RR_OVERRIDE_RENDER_PRESET_SELECTION_ID);
                        if (rrSetting != null && rrSetting.CurrentValue is uint)
                            dlssdPreset = (uint)rrSetting.CurrentValue;

                        Console.WriteLine(string.Format("{{\"success\": true, \"dlssPreset\": {0}, \"dlssdPreset\": {1}, \"foundProfile\": true}}", dlssPreset, dlssdPreset));
                    }
                    return 0;
                }
                else if (command == "set")
                {
                    if (args.Length < 4)
                    {
                        Console.WriteLine("{\"success\": false, \"error\": \"Missing tech or preset arguments\"}");
                        return 1;
                    }

                    string tech = args[2].ToLowerInvariant();
                    uint presetValue;
                    if (!uint.TryParse(args[3], out presetValue))
                    {
                        Console.WriteLine("{\"success\": false, \"error\": \"Invalid preset value\"}");
                        return 1;
                    }

                    if (profile == null)
                    {
                        Console.WriteLine("{\"success\": false, \"error\": \"Driver profile not found\"}");
                        return 1;
                    }

                    if (tech == "dlss")
                        profile.SetSetting(NGX_DLSS_SR_OVERRIDE_RENDER_PRESET_SELECTION_ID, presetValue);
                    else if (tech == "dlssd")
                        profile.SetSetting(NGX_DLSS_RR_OVERRIDE_RENDER_PRESET_SELECTION_ID, presetValue);
                    else
                    {
                        Console.WriteLine(string.Format("{{\"success\": false, \"error\": \"Unknown tech {0}\"}}", tech));
                        return 1;
                    }

                    session.Save();
                    Console.WriteLine("{\"success\": true}");
                    return 0;
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(string.Format("{{\"success\": false, \"error\": \"{0}\"}}", ex.Message.Replace("\"", "'").Replace("\n", " ").Replace("\r", "")));
                return 1;
            }
            return 0;
        }

        static DriverSettingsProfile FindProfile(DriverSettingsSession session, string exeName)
        {
            foreach (var profile in session.Profiles)
            {
                foreach (var app in profile.Applications)
                {
                    if (string.Equals(app.ApplicationName, exeName, StringComparison.OrdinalIgnoreCase))
                    {
                        return profile;
                    }
                }
            }
            return null;
        }
    }
}
