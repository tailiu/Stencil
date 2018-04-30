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

        # if @user != nil
        #     @new_tweet = Tweet.new(content: params[:tweet], reply_to: params[:reply_to], user_id: @user["id"])
        #     if @new_tweet.valid?
        #         @new_tweet.save
        #         @result["success"] = true
        #         @result["tweet"] = @new_tweet
        #     else
        #         @result["success"] = false
        #         @result["error"]["message"] = "Couldn't create new tweet. Check params."
        #     end
        # else
        #     @result["success"] = false
        #     @result["error"]["message"] = "User doesn't exist!"
        #     @result["error"]["params"] = params
        # end

        @result[:session] = session
        @result[:session_id] = request.session_options[:id]
        render json: {result: @result}
    end

    def fetchall
        @result = {
            "success" => true,
            "error" => {
            },
        }

        @tweets = Tweet.all
        @result["tweets"] = @tweets
        
        render json: {result: @result}
    end


    def fetchallbyuser
        
        @user = User.find_by_id(params[:user_id])
        
        @result = {
            "success" => true,
            "error" => {
            },
        }

        if @user == nil
            @result["error"]["message"] = "User doesn't exist!"
            @result["success"] = false
        else
            @result["success"] = true
            @tweets = Tweet.where(:user_id => params[:user_id])
            @result["tweets"] = @tweets
        end

        
        render json: {result: @result}
    end


    def like
    end

    def retweet
    end

end
