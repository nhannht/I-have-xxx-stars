package dev.nhannht;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.LambdaLogger;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.kohsuke.github.GHEventPayload;
import org.kohsuke.github.GitHub;

import java.io.IOException;
import java.io.StringReader;
import java.util.HashMap;
import java.util.Map;

public class XXXStars implements RequestHandler<APIGatewayProxyRequestEvent, APIGatewayProxyResponseEvent> {


    @Override
    public APIGatewayProxyResponseEvent handleRequest(APIGatewayProxyRequestEvent apiGatewayProxyRequestEvent, Context context) {
        LambdaLogger logger = context.getLogger();
        Map<String, String> headers = new HashMap<>();
        headers.put("Content-Type", "application/json");
        headers.put("X-Custom-Header", "application/json");

        APIGatewayProxyResponseEvent response = new APIGatewayProxyResponseEvent()
                .withHeaders(headers);
        var reader = new StringReader(apiGatewayProxyRequestEvent.getBody());
        SecretManager secretManager = new SecretManager();
        String secret = secretManager.getSecret();
        var mapper = new ObjectMapper();
        try {
            Secret secretPOJO = mapper.readValue(secret, Secret.class);
            String oauth = secretPOJO.getGithub_oauth_repo1();
            GitHub githubClient = GitHub.connectUsingOAuth(oauth);

            var payload = githubClient.parseEventPayload(reader,
                    GHEventPayload.Star.class);
            var repo = payload.getRepository();
            var starGazers = repo.getStargazersCount();

            repo.renameTo(String.format("I-have-%s-stars", starGazers));
            repo.setDescription(String.format("I have %s stars",starGazers));
            return response.withStatusCode(200)
                    .withBody(repo.getName());

        } catch (IOException e) {
            throw new RuntimeException(e);
        }

    }


}