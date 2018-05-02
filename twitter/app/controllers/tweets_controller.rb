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
            "params" => params,
            "success" => false,
            "error" => {
            },
        }

        if params[:content].nil? || params[:user_id].nil?
            @result["success"] = false
            @result["error"]["message"] = "Incomplete params!"
        elsif params[:content].empty?
            @result["success"] = false
            @result["error"]["message"] = "Tweet can't be empty!"
        else
            @user = User.find_by_id(params[:user_id])
            if @user != nil
                @new_tweet = Tweet.new(content: params[:content], user_id: params[:user_id], reply_to_id: nil)
                if @new_tweet.valid?
                    @new_tweet.save
                    @result["success"] = true
                    @result["tweet"] = @new_tweet
                else
                    puts @new_tweet.errors.messages
                    @result["success"] = false
                    @result["error"]["message"] = @new_tweet.errors.messages
                end
            else
                @result["success"] = false
                @result["error"]["message"] = "User doesn't exist!"
            end
        end

        render json: {result: @result}
    end


    def fetchUserTweets

        result = {
            "params" => params,
            "success" => false,
            "error" => {
            },
        }

        if params[:user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif !User.where(id: params[:user_id]).exists?
            result["success"] = false
            result["error"]["message"] = "User doesn't exist!"
        else
            user = User.find(params[:user_id])
            # @user = User.where(id: params[:user_id]).first
            tweets = user.tweets
            # @tweets = Tweet.where(user_id: params[:user_id]).order('created_at DESC')
            result["success"] = true
            tweets_set = []
            for tweet in tweets do
                tweets_set.push({"tweet": tweet, "creator": tweet.user, "likes": tweet.likes.count, "retweets": tweet.retweets.count, "replies": tweet.replies.count})
            end
            result["tweets"] = tweets_set
        end

        render json: {result: result}

    end


    def mainPageTweets
        result = {
            "params" => params,
            "success" => false,
            "error" => {
            },
        }

        if params[:user_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif !User.where(id: params[:user_id]).exists?
            result["success"] = false
            result["error"]["message"] = "User doesn't exist!"
        else
            tweets_set = []
            user = User.find(params[:user_id])
            following_users = UserAction.where(from_user_id: user.id, action_type: "follow").pluck(:to_user_id)
            allusers = following_users + [user.id]
            alltweets = Tweet.where(:user_id => allusers).order('created_at DESC')
            for tweet in alltweets do
                tweets_set.push({"tweet": tweet, "creator": tweet.user,  "likes": tweet.likes.count, "retweets": tweet.retweets.count, "replies": tweet.replies.count})
            end
            result["success"] = true
            result["tweets"] = tweets_set
        end

        render json: {result: result}
    end

    def getTweet
        result = {
            "params" => params,
            "success" => false,
            "error" => {
            },
        }
        tweet = Tweet.find_by_id(params[:tweet_id])
        if !tweet.nil?
            # tweets = Tweet.where(reply_to_id: tweet.id)
            tweets = tweet.replies
            tweets_set = [{"tweet": tweet, "creator": tweet.user,  "likes": tweet.likes.count, "retweets": tweet.retweets.count, "replies": tweet.replies.count}]
            for tweet in tweets do
                tweets_set.push({"tweet": tweet, "creator": tweet.user,  "likes": tweet.likes.count, "retweets": tweet.retweets.count, "replies": tweet.replies.count})
            end
            result["replies"] = tweets_set
            result["success"] = true
        else
            result["success"] = false
            result["error"]["message"] = "Tweet doesn't exist!"
        end
        render json: {result: result}
    end


    def like
        result = {
            "params" => params,
            "success" => false,
            "error" => {
            },
        }

        if params[:tweet_id].nil? || params[:user_id].nil? || params[:like].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.find_by_id(params[:user_id])
            tweet = Tweet.find_by_id(params[:tweet_id])
            if !user.nil? && !tweet.nil?
                # like = Like.new(tweet_id: params[:tweet_id], user_id: params[:user_id])
                # like = tweet.likes.create(tweet_id: params[:tweet_id], user_id: params[:user_id])
                like = tweet.likes.find_or_create_by(tweet_id: params[:tweet_id], user_id: params[:user_id])
                if like.valid?
                    result["success"] = true
                    if params[:like] === "true"
                        like.save
                        result["like"] = like
                    else
                        Like.destroy(like.id)
                    end
                else
                    puts like.errors.messages
                    result["success"] = false
                    result["error"]["message"] = like.errors.messages
                end
            else
                result["success"] = false
                result["error"]["message"] = "User/Tweet don't exist!"
            end
        end

        render json: {result: result}
    end

    def retweet
        result = {
            "params" => params,
            "success" => false,
            "error" => {
            },
        }

        if params[:tweet_id].nil? || params[:user_id].nil? || params[:retweet].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.find_by_id(params[:user_id])
            tweet = Tweet.find_by_id(params[:tweet_id])
            if !user.nil? && !tweet.nil?
                # like = Like.new(tweet_id: params[:tweet_id], user_id: params[:user_id])
                # like = tweet.likes.create(tweet_id: params[:tweet_id], user_id: params[:user_id])
                retweet = tweet.retweets.find_or_create_by(tweet_id: params[:tweet_id], user_id: params[:user_id])
                if retweet.valid?
                    result["success"] = true
                    if params[:retweet] === "true"
                        retweet.save
                        result["retweet"] = retweet
                    else
                        Retweet.destroy(retweet.id)
                    end
                else
                    puts retweet.errors.messages
                    result["success"] = false
                    result["error"]["message"] = retweet.errors.messages
                end
            else
                result["success"] = false
                result["error"]["message"] = "User/Tweet don't exist!"
            end
        end

        render json: {result: result}
    end

    def stats
        result = {
            "params" => params,
            "success" => false,
            "error" => {
            },
        }

        if params[:tweet_id].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            result["success"] = true
            tweet = Tweet.find_by_id(params[:tweet_id])
            result["likes"] = tweet.likes.pluck(:user_id)
            result["retweets"] = tweet.retweets.pluck(:user_id)
            result["replies"] = tweet.replies.pluck(:user_id)
            # result["replies_users"] = tweet.replies
        end
        
        render json: {result: result}
    end

    def reply 
    end

end
