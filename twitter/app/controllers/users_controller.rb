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

            session[:user_id] = @new_user.id
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
<<<<<<< HEAD
        @credential = Credential.find_by(email: params[:email], password: params[:password])
=======

        reset_session

        @credentials = Credential.find_by(email: params[:email], password: params[:password])
>>>>>>> 959e0342a9bfe79b7c195812a57ca285aab381a4

        @result = {
            # params: params,
            "success" => false,
            "error" => {
            }
        }

<<<<<<< HEAD
        if @credential != nil
            @result["success"] = true
            @result["user"]  = @credential.user

            session[:user_id] = @credential.user.id
=======
        if @credentials != nil
            @result["success"] = true
            @result["user"]  = @credentials.user
            session[@credentials.user.id]
            @result["session_id"]  = session.id

>>>>>>> 959e0342a9bfe79b7c195812a57ca285aab381a4
        else
            @result["success"] = false
            @result["error"]["message"] = "Invalid credentials!"
        end

        render json: {result: @result}
    end

    def logout
        session.clear
    end


end
