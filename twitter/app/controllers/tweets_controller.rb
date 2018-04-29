class TweetsController < ApplicationController
    def index
        if params[:type] == "tweet_num"
            @user = User.find(params[:id])
            @tweet_num = @user.tweets.size

            @result = {
                # params: params,
                "tweet_num" => @tweet_num,
            } 

            render json: {result: @result}
        end
    end

    def new
        @result = {
            # params: params,
            "success" => false,
            "error" => {
            },
            "user" => session[:user]
        }

        puts session

        # @user = User.find_by(handle: params[:handle])
        # @result["user"] = @user

        @user = session[:user]

        if @user != nil
            @new_tweet = Tweet.new(content: params[:tweet], reply_to: params[:reply_to], user_id: @user["id"])
            if @new_tweet.valid?
                @new_tweet.save
                @result["success"] = true
                @result["tweet"] = @new_tweet
            else
                @result["success"] = false
                @result["error"]["message"] = "Couldn't create new tweet. Check params."
            end
        else
            @result["success"] = false
            @result["error"]["message"] = "User doesn't exist!"
            @result["error"]["params"] = params
        end

        @result[:session] = session
        render json: {result: @result}
    end

end
