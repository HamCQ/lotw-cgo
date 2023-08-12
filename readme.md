Using tqsllib
-------------
A minimal set of calls to permit an application to sign a log is the following.
Of course, error checking should be performed for each call.

        tqsl_getStationLocation(&loc, location_name);

        tqsl_getLocationCallSign(loc, callsign, sizeof callsign);
        tqsl_getLocationDXCCEntity(loc, &dxcc);
        tqsl_selectCertificates(&certlist, &ncerts, callsign, dxcc);

        tqsl_beginADIFConverter(&conv, input_file, certlist, ncerts, loc);
        tqsl_setConverterAllowDuplicates(conv, false);

        tqsl_setConverterAppName(conv, "myAppName");
            (tell tqsllib the name of your application)

        while (cp = tqsl_getConverterGABBI(conv) != 0)
            write the string pointed to by "cp" to your file

        tqsl_converterCommit(conv);
        tqsl_endConverter(&conv);
        tqsl_endStationLocationCapture(&loc);

The tq8 files created by tqsl are compressed using zlib functions. You can
also submit uncompressed files using a .tq7 extension.