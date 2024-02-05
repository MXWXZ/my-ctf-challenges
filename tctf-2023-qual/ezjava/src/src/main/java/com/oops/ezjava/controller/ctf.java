package com.oops.ezjava.controller;

import com.googlecode.aviator.AviatorEvaluatorInstance;
import com.googlecode.aviator.Feature;
import com.googlecode.aviator.Options;
import com.googlecode.aviator.script.AviatorScriptEngine;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.script.ScriptEngine;
import javax.script.ScriptEngineManager;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.Base64;
import java.util.HashSet;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

@RestController
public class ctf {
    @PostMapping("/eval")
    public String eval(@RequestParam(value = "url", defaultValue = "http://localhost:8080/aviatorscript") String url) {
        Pattern pattern = Pattern.compile("^[a-z]+[a-zA-Z0-9:/]+$");
        if (!pattern.matcher(url).find()) {
            return "No";
        }
        String low = url.toLowerCase();
        if (low.contains("file") || low.contains("dict")) {
            return "No";
        }
        if (low.startsWith("http")) {
            if (!low.startsWith("http://localhost:8080/")) {
                return "http host not allowed";
            }
        }

        String[] cmd = {"./new_curl", "-s", url};
        try {
            Process pr = Runtime.getRuntime().exec(cmd);
            String result = new BufferedReader(new InputStreamReader(pr.getInputStream()))
                    .lines()
                    .collect(Collectors.joining("\n"));
            if (!result.isEmpty()) {
                System.out.printf("url: %s result: %s\n", url, Base64.getEncoder().encodeToString(result.getBytes()));
                ScriptEngineManager sem = new ScriptEngineManager();
                ScriptEngine engine = sem.getEngineByName("AviatorScript");
                AviatorEvaluatorInstance instance = ((AviatorScriptEngine) engine).getEngine();
                instance.setOption(Options.FEATURE_SET, Feature.asSet());
                instance.setOption(Options.ALLOWED_CLASS_SET, new HashSet());
                return engine.eval(result).toString();
            }
        } catch (Exception e) {
            return "error";
        }
        return "no response";
    }

    @GetMapping("/aviatorscript")
    public String script() {
        return "114514+1919810";
    }
}
