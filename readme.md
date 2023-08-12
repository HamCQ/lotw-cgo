# lotw-go

本项目使用 Go 调用 tqsllib.dll 动态库的一些操作  

tqsllib 源码同步放在 src 文件夹中，更完整的源码可以前往参考材料中提到的 sourceforge 中下载    

使用前需要先查看 build.sh 中提到的 32 位环境配置  

动态库的存储文件，会生成在 C:\Users\Emin\AppData\Roaming\TrustedQSL 中保存，Emin 为对应的电脑中的用户  



# 参考材料
* http://www.arrl.org/files/file/LoTW%20Instructions/TQSL%202_0%202013/Step%205%20Defining%20A%20Station%20Locatione_2013.pdf
* https://www.hellocq.net/forum/read.php?tid=371255
* https://www.rickmurphy.net/lotw/doxygen/index.html
* http://www.arrl.org/developer
* https://sourceforge.net/projects/trustedqsl/files/tqsllib/

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