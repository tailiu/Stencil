class UsersController < ApplicationController

    def getUserInfo

        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }

        if params[:user_id].nil? || params[:req_token].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif session[:req_token].to_s != params[:req_token]
            result["success"] = false
            result["error"]["message"] = "Invalid token!"
        else
            user = User.find_by_id(params[:user_id])
            if !user.nil?
                result["user"] = user
                result["email"] = user.credential.email
                result["user_stats"] = {
                    "tweets" => user.tweets.size,
                    "followers" => UserAction.where(to_user_id: user.id, action_type: "follow").count,
                    "following" => UserAction.where(from_user_id: user.id, action_type: "follow").count
                }
                result["success"] = true
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def userInfo
    end

    def updateBio
        result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }
        if params[:user_id].nil? || params[:bio].nil? || params[:req_token].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif session[:req_token].to_s != params[:req_token]
            result["success"] = false
            result["error"]["message"] = "Invalid token!"
        else
            user = User.find_by_id(params[:user_id])
            if user != nil
                result["user"] = user
                result["success"] = true
                user.bio = params[:bio]
                user.save
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def updatePhoto
    end

    def updateEmail
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        if params[:user_id].nil? || params[:email].nil? || params[:req_token].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif session[:req_token].to_s != params[:req_token]
            result["success"] = false
            result["error"]["message"] = "Invalid token!"
        else
            user = Credential.where(user_id: params[:user_id]).first
            if user != nil
                result["success"] = true
                user.email = params[:email]
                user.save
                result["user"] = user
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def updateHandle
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        if params[:user_id].nil? || params[:handle].nil? || params[:req_token].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif session[:req_token].to_s != params[:req_token]
            result["success"] = false
            result["error"]["message"] = "Invalid token!"
        else
            user = User.find_by_id(params[:user_id])
            if user != nil
                result["success"] = true
                user.handle = params[:handle]
                user.save
                result["user"] = user
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def updatePassword
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        if params[:user_id].nil? || params[:password].nil? || params[:req_token].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif session[:req_token].to_s != params[:req_token]
            result["success"] = false
            result["error"]["message"] = "Invalid token!"
        else
            user = Credential.where(user_id: params[:user_id]).first
            if user != nil
                result["success"] = true
                user.password = params[:password]
                user.save
                result["user"] = user
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def updateProtected
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        if params[:user_id].nil? || params[:protected].nil? || params[:req_token].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        elsif session[:req_token].to_s != params[:req_token]
            result["success"] = false
            result["error"]["message"] = "Invalid token!"
        else
            user = User.find_by_id(params[:user_id])
            if user != nil
                result["success"] = true
                if params[:protected] === "true"
                    user.protected = true
                else
                    user.protected = false
                end
                user.save
                result["user"] = user
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def search
        result = {
            "params": params,
            "success" => false,
            "error" => {
            }
        }
        if params[:query].nil? || params[:query].empty?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            result["success"] = true
            query = params[:query]
            result["users"] = User.where("lower(handle) like ?", "%#{query.downcase}%").or(User.where("lower(name) like ?", "%#{query.downcase}%"))
        end
        render json: {result: result}
    end

end
