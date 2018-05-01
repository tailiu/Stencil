class UsersController < ApplicationController
    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            @result["success"] = true
            @result["user"] = @new_user

            session[:current_user_id] = @new_user.id
        else 
            puts @new_user.errors.messages
            puts @new_credential.errors.messages

            if @new_user.errors.messages.key?(:handle) && @new_credential.errors.messages.key?(:email)
                @result["error"]["message"] = "Handle and Email already exist!"
            elsif @new_user.errors.messages.key?(:handle)
                @result["error"]["message"] = "Handle already exists!"
            elsif @new_credential.errors.messages.key?(:email)
                @result["error"]["message"] = "Email already exists!"
            else
                @result["error"]["message"] = "Invalid Credentials!"
            end

            @result["success"] = false
        end

        render json: {result: @result}
    end

    def verify

        @credentials = Credential.find_by(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

        if @credentials != nil
            @result["success"] = true
            @result["user"]  = @credentials.user
            
            session[:current_user_id] = nil
            session[:current_user_id] = @credentials.user.id
            @result["session_id"]  = session.id

        else
            @result["success"] = false
            @result["error"]["message"] = "Invalid credentials!"
        end

        @result[:session] = session
        render json: {result: @result}
    end

    def logout
        session[:current_user_id] = nil
        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        render json: {result: @result}
    end

    def getUserInfo
        @user = User.find_by_id(params[:user_id])
        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        if @user != nil
            @result["user"] = @user
            @result["user_stats"] = {
                "tweets" => @user.tweets.size,
                "followers" => UserAction.where(to_user_id: @user.id, action_type: "follow").count,
                "following" => UserAction.where(from_user_id: @user.id, action_type: "follow").count
            }
            @result["success"] = true
        else
            @result["success"] = false
            @result["error"]["message"] = "User doesn't exist!"
        end
        
        render json: {result: @result}
    end

    def userInfo
    end

    def updateBio
        result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        if params[:user_id].nil? || params[:bio].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
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
        if params[:user_id].nil? || params[:email].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = Credential.where(user_id: params[:user_id]).first
            if user != nil
                result["user"] = user
                result["success"] = true
                user.email = params[:email]
                user.save
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
        if params[:user_id].nil? || params[:handle].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = User.find_by_id(params[:user_id])
            if user != nil
                result["user"] = user
                result["success"] = true
                user.handle = params[:handle]
                user.save
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
        if params[:user_id].nil? || params[:password].nil?
            result["success"] = false
            result["error"]["message"] = "Incomplete params!"
        else
            user = Credential.where(user_id: params[:user_id]).first
            if user != nil
                result["user"] = user
                result["success"] = true
                user.password = params[:password]
                user.save
            else
                result["success"] = false
                result["error"]["message"] = "User doesn't exist!"
            end
        end
        
        render json: {result: result}
    end

    def markAsProtected
    end

end
